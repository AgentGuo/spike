/*
@author: panfengguo
@since: 2024/11/17
@desc: desc
*/
package reqscheduler

import (
	"context"
	"fmt"
	"github.com/AgentGuo/spike/api"
	"github.com/AgentGuo/spike/cmd/server/config"
	"github.com/AgentGuo/spike/pkg/logger"
	"github.com/AgentGuo/spike/pkg/storage"
	"github.com/AgentGuo/spike/pkg/storage/model"
	"github.com/AgentGuo/spike/pkg/worker"
	"github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sort"
	"time"
)

type Request struct {
	FunctionName    string
	RequestID       uint64
	ReqPayload      string
	RequiredCpu     int32
	RequiredMemory  int32
	RespPayloadChan chan Response
}

type Response struct {
	ResponsePayload string
	err             error
}

type ReqScheduler struct {
	mysql     *storage.Mysql
	logger    *logrus.Logger
	reqQueue  *ReqQueue
	flake     *sonyflake.Sonyflake
	triggerCh chan struct{}
}

func NewReqScheduler() *ReqScheduler {
	mysqlClient := storage.NewMysql()
	r := &ReqScheduler{
		mysql:     mysqlClient,
		logger:    logger.GetLogger(),
		reqQueue:  NewReqQueue(),
		flake:     sonyflake.NewSonyflake(sonyflake.Settings{}),
		triggerCh: make(chan struct{}),
	}
	go r.ScheduleRoutine()
	return r
}

func (r *ReqScheduler) ScheduleRoutine() {
	// 创建一个定时器，每秒触发一次
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			break
		case <-r.triggerCh:
			break
		}
		r.Schedule()
	}
}

func (r *ReqScheduler) Schedule() {
	// step1: get request from queue
	req := r.reqQueue.Peek()
	if req == nil {
		return
	}

	// step2: get function instance
	funcInstances, err := r.mysql.GetFuncInstanceByCondition(map[string]interface{}{
		"function_name": req.FunctionName,
		"last_status":   "RUNNING",
	})
	if err != nil {
		r.logger.Errorf("get function instance failed, %v", err)
		return
	}

	// step4: get processing request
	reqScheduleInfo, err := r.mysql.GetReqScheduleInfoByFunctionName(req.FunctionName)
	if err != nil {
		r.logger.Errorf("get request schedule info failed, %v", err)
		return
	}

	// step5: rank function instance
	type instanceStat struct {
		awsServiceName  string
		ipv4            string
		cpu             int32
		memory          int32
		cpuUsed         int32
		memoryUsed      int32
		cpuUsageRate    float64
		memoryUsageRate float64
		avgUsageRate    float64
	}
	insStatMap := make(map[string]*instanceStat)
	for _, instance := range funcInstances {
		insStatMap[instance.AwsServiceName] = &instanceStat{
			awsServiceName: instance.AwsServiceName,
			ipv4:           instance.PublicIpv4,
			cpu:            instance.Cpu,
			memory:         instance.Memory,
			cpuUsed:        0,
			memoryUsed:     0,
		}
	}
	for _, reqInfo := range reqScheduleInfo {
		insStatMap[reqInfo.PlacedAwsServiceName].cpuUsed = reqInfo.RequiredCpu
		insStatMap[reqInfo.PlacedAwsServiceName].memoryUsed = reqInfo.RequiredMemory
	}
	insStatList := make([]*instanceStat, 0, len(insStatMap))
	for _, v := range insStatMap {
		v.cpuUsageRate = float64(v.cpuUsed) / float64(v.cpu)
		v.memoryUsageRate = float64(v.memoryUsed) / float64(v.memory)
		v.avgUsageRate = (v.cpuUsageRate + v.avgUsageRate) / 2
		insStatList = append(insStatList, v)
	}

	sort.Slice(insStatList, func(i, j int) bool {
		if insStatList[i].cpu != insStatList[j].cpu {
			return insStatList[i].cpu < insStatList[j].cpu
		} else if insStatList[i].memory != insStatList[j].memory {
			return insStatList[i].memory < insStatList[j].memory
		} else {
			return insStatList[i].avgUsageRate > insStatList[j].avgUsageRate
		}
	})

	// step6: chose function instance to send request
	var chosenInsIpv4, choseAwsServiceName string
	for _, insStat := range insStatList {
		if insStat.cpuUsed+req.RequiredCpu <= insStat.cpu && insStat.memoryUsed+req.RequiredMemory <= insStat.memory {
			chosenInsIpv4 = insStat.ipv4
			choseAwsServiceName = insStat.awsServiceName
			break
		}
	}
	if chosenInsIpv4 == "" {
		r.logger.Errorf("no available instance to handle request")
		return
	}
	newReqScheduleInfo := &model.ReqScheduleInfo{
		ReqId:                req.RequestID,
		FunctionName:         req.FunctionName,
		PlacedAwsServiceName: choseAwsServiceName,
		RequiredCpu:          req.RequiredCpu,
		RequiredMemory:       req.RequiredMemory,
	}
	err = r.mysql.UpdateReqScheduleInfo(newReqScheduleInfo)
	if err != nil {
		r.logger.Errorf("update req schedule info failed, %v", err)
		return
	}
	r.reqQueue.Pop()
	r.logger.Infof("schedule request %d to node: %s(%s)", req.RequestID, choseAwsServiceName, chosenInsIpv4)
	go r.CallInstanceFunctionRoutine(req, chosenInsIpv4)
}

func (r *ReqScheduler) SubmitRequest(req *api.CallFunctionRequest, respChan chan Response) error {
	// step1: construct request
	reqID, err := r.flake.NextID()
	if err != nil {
		return err
	}
	request := &Request{
		FunctionName:    req.GetFunctionName(),
		RequestID:       reqID,
		ReqPayload:      req.Payload,
		RequiredCpu:     req.Cpu,
		RequiredMemory:  req.Memory,
		RespPayloadChan: respChan,
	}

	// step2: submit into request queue
	r.reqQueue.Push(request)
	r.triggerCh <- struct{}{}
	return nil
}

// CallFunction 对外暴露的函数调用接口
func (r *ReqScheduler) CallFunction(req *api.CallFunctionRequest) (*api.CallFunctionResponse, error) {
	respChan := make(chan Response)
	err := r.SubmitRequest(req, respChan)
	if err != nil {
		return nil, err
	}
	resp := <-respChan
	if resp.err != nil {
		return nil, resp.err
	}
	return &api.CallFunctionResponse{ErrorCode: 0, Payload: resp.ResponsePayload}, nil
}

// CallInstanceFunctionRoutine 调用实例函数的协程
func (r *ReqScheduler) CallInstanceFunctionRoutine(req *Request, instanceIpv4 string) {
	respPayload, err := r.CallInstanceFunction(req.ReqPayload, req.RequestID, instanceIpv4)
	resp := Response{
		ResponsePayload: respPayload,
		err:             err,
	}
	req.RespPayloadChan <- resp
	err = r.mysql.DeleteReqScheduleInfo(req.RequestID)
	if err != nil {
		r.logger.Errorf("delete req schedule info failed, %v", err)
	}
}

// CallInstanceFunction 调用实例函数
func (r *ReqScheduler) CallInstanceFunction(reqPayload string, reqID uint64, instanceIpv4 string) (string, error) {
	// TODO: 这里可以做连接复用
	r.logger.Infof("call function %s", instanceIpv4)
	conn, err := grpc.NewClient(fmt.Sprintf("%s:50052", instanceIpv4), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close() // 确保连接关闭
	workerServiceClient := worker.NewSpikeWorkerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GetConfig().DispatchTimeout)*time.Second)
	defer cancel()
	funcServiceResp, err := workerServiceClient.CallWorkerFunction(ctx, &worker.CallWorkerFunctionReq{
		Payload:   reqPayload,
		RequestId: reqID,
	})
	if err != nil {
		return "", err
	}
	return funcServiceResp.Payload, nil
}

//func (f *ReqScheduler) CallFunction(req *api.CallFunctionRequest) (*api.CallFunctionResponse, error) {
//	funcInstances, err := f.mysql.GetFuncInstanceByCondition(map[string]interface{}{
//		"function_name": req.FunctionName,
//		"cpu":           req.Cpu,
//		"memory":        req.Memory,
//		"last_status":   "RUNNING",
//	})
//	if err != nil {
//		return nil, err
//	} else if len(funcInstances) == 0 {
//		return nil, fmt.Errorf("not such a function or function is not ready now")
//	}
//	task := funcInstances[rand.Intn(len(funcInstances))]
//	f.logger.Infof("call function %s", task.PublicIpv4)
//	conn, err := grpc.NewClient(fmt.Sprintf("%s:50052", task.PublicIpv4), grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		return nil, err
//	}
//	workerServiceClient := worker.NewSpikeWorkerServiceClient(conn)
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GetConfig().DispatchTimeout)*time.Second)
//	defer cancel()
//	funcServiceResp, err := workerServiceClient.CallWorkerFunction(ctx, &worker.CallWorkerFunctionReq{
//		Payload:   req.Payload,
//		RequestId: strconv.FormatInt(time.Now().Unix(), 10),
//	})
//	if err != nil {
//		return nil, err
//	}
//	return &api.CallFunctionResponse{ErrorCode: 0, Payload: funcServiceResp.Payload}, nil
//}
