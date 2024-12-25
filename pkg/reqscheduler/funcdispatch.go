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
	"github.com/AgentGuo/spike/pkg/worker"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"strconv"
	"time"
)

// Request 表示一个请求
type Request struct {
	ID      int64
	Payload string
}

type ReqScheduler struct {
	mysqlClient *storage.Mysql
	logger      *logrus.Logger
}

func NewReqScheduler() *ReqScheduler {
	mysqlClient := storage.NewMysql()
	return &ReqScheduler{
		mysqlClient: mysqlClient,
		logger:      logger.GetLogger(),
	}
}

func (f *ReqScheduler) AddRequest(req *api.CallFunctionRequest) {
	// TODO
}

func (f *ReqScheduler) CallFunction(req *api.CallFunctionRequest) (*api.CallFunctionResponse, error) {
	taskData, err := f.mysqlClient.GetFuncInstanceByCondition(map[string]interface{}{
		"function_name": req.FunctionName,
		"cpu":           req.Cpu,
		"memory":        req.Memory,
		"last_status":   "RUNNING",
	})
	if err != nil {
		return nil, err
	} else if len(taskData) == 0 {
		return nil, fmt.Errorf("not such a function or function is not ready now")
	}
	task := taskData[rand.Intn(len(taskData))]
	f.logger.Infof("call function %s", task.PublicIpv4)
	conn, err := grpc.NewClient(fmt.Sprintf("%s:50052", task.PublicIpv4), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	workerServiceClient := worker.NewSpikeWorkerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GetConfig().DispatchTimeout)*time.Second)
	defer cancel()
	funcServiceResp, err := workerServiceClient.CallWorkerFunction(ctx, &worker.CallWorkerFunctionReq{
		Payload:   req.Payload,
		RequestId: strconv.FormatInt(time.Now().Unix(), 10),
	})
	if err != nil {
		return nil, err
	}
	return &api.CallFunctionResponse{ErrorCode: 0, Payload: funcServiceResp.Payload}, nil
}
