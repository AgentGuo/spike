/*
@author: panfengguo
@since: 2024/11/6
@decs: decs
*/
package funcmanager

import (
	"container/list"
	"fmt"
	"github.com/AgentGuo/spike/api"
	"github.com/AgentGuo/spike/cmd/server/config"
	"github.com/AgentGuo/spike/pkg/logger"
	"github.com/AgentGuo/spike/pkg/reqscheduler"
	"github.com/AgentGuo/spike/pkg/storage"
	"github.com/AgentGuo/spike/pkg/storage/model"
	"github.com/AgentGuo/spike/pkg/utils"
	"github.com/sirupsen/logrus"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"
)

type FuncManager struct {
	awsClient *AwsClient
	mysql     *storage.Mysql
	logger    *logrus.Logger
	reqQueue  *reqscheduler.ReqQueue
}

var (
	funcManager *FuncManager
	once        sync.Once
)

func NewFuncManager() *FuncManager {
	once.Do(func() {
		funcManager = &FuncManager{
			awsClient: NewAwsClient(),
			mysql:     storage.NewMysql(),
			logger:    logger.GetLogger(),
			reqQueue:  reqscheduler.NewReqQueue(),
		}
		go funcManager.FunctionAutoScale()
	})
	return funcManager
}

func (f *FuncManager) CreateFunction(req *api.CreateFunctionRequest) error {
	// step1: check if function has been created
	if hasCreated, err := funcManager.mysql.HasFuncMetaDataByFunctionName(req.FunctionName); err != nil || hasCreated {
		if hasCreated {
			return fmt.Errorf("function has been created")
		} else {
			return err
		}
	}

	// step2: create function definition
	var resourceSpecList []model.ResourceSpec
	for _, res := range req.Resources {
		var family string
		var revision int32
		if awsTaskDef, err := f.mysql.GetAwsTaskDefByFuncCpuMenImg(req.FunctionName, res.Cpu, res.Memory, req.ImageUrl); err == nil && awsTaskDef != nil {
			family = awsTaskDef.TaskFamily
			revision = awsTaskDef.TaskRevision
		} else {
			family, revision, err = f.awsClient.RegTaskDef(req.FunctionName, res.Cpu, res.Memory, req.ImageUrl)
			if err != nil {
				return err
			}
			updateTaskDef := &model.AwsTaskDef{
				TaskFamily:   family,
				TaskRevision: revision,
				FunctionName: req.FunctionName,
				Cpu:          res.Cpu,
				Memory:       res.Memory,
				ImageUrl:     req.ImageUrl,
			}
			err = f.mysql.UpdateAwsTaskDef(updateTaskDef)
			if err != nil {
				f.logger.Errorf("update aws task def failed, err: %v", err)
			}
		}
		resourceSpecList = append(resourceSpecList, model.ResourceSpec{
			Cpu:        res.Cpu,
			Memory:     res.Memory,
			MinReplica: res.MinReplica,
			MaxReplica: res.MaxReplica,
			Family:     family,
			Revision:   revision,
		})
	}
	err := funcManager.mysql.CreateFuncMetaData(&model.FuncMetaData{
		FunctionName: req.FunctionName,
		ImageUrl:     req.ImageUrl,
		ResSpecList:  resourceSpecList,
	})
	if err != nil {
		f.logger.Errorf("create func meta data failed, err: %v", err)
		return err
	}

	// step3: create function instance
	var funcInstances []model.FuncInstance
	for _, res := range resourceSpecList {
		// TODO: test only
		awsServiceNames, err := f.awsClient.BatchCreateInstance(res.Family, res.Revision, Fargate, res.MinReplica)
		if err != nil {
			f.logger.Errorf("create ecs failed, err: %v", err)
			return err
		}
		for _, awsServiceName := range awsServiceNames {
			funcInstances = append(funcInstances, model.FuncInstance{
				AwsServiceName: awsServiceName,
				FunctionName:   req.FunctionName,
				Cpu:            res.Cpu,
				Memory:         res.Memory,
				AwsFamily:      res.Family,
				AwsRevision:    res.Revision,
				LastStatus:     "NOT_CREATE",
				//LaunchType:     int32(EC2),
				LaunchType: int32(Fargate), // TODO: test only
			})
		}
	}
	if err := f.mysql.UpdateFuncInstanceBatch(funcInstances); err != nil {
		f.logger.Errorf("UpdateFuncInstanceBatch failed, err: %v", err)
		return err
	}
	go f.UpdateFunctionStatus(req.FunctionName)
	return nil
}

func (f *FuncManager) ScaleFunction(req *api.ScaleFunctionRequest) error {
	// step1: 获取当前函数信息
	metaData, err := f.mysql.GetFuncMetaDataByFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("scale function failed, get func meta data failed, functionName: %s, err: %v", req.FunctionName, err)
		return err
	}
	currentInsList, err := f.mysql.GetFuncInstanceByCondition(map[string]interface{}{"function_name": req.FunctionName,
		"cpu":    req.Cpu,
		"memory": req.Memory,
	})
	if err != nil {
		f.logger.Errorf("get func instance failed, functionName: %s, err: %v", req.FunctionName, err)
		return err
	}
	awsTaskDef, err := f.mysql.GetAwsTaskDefByFuncCpuMenImg(req.FunctionName, req.Cpu, req.Memory, metaData.ImageUrl)
	if err != nil {
		f.logger.Errorf("scale function failed, get aws task def failed, functionName: %s, err: %v", req.FunctionName, err)
		return err
	}

	// step2: 检查是否超出最大实例数
	maxReplica, minReplica := int32(-1), int32(-1)
	for _, res := range metaData.ResSpecList {
		if res.Cpu == req.Cpu && res.Memory == req.Memory {
			maxReplica, minReplica = res.MaxReplica, res.MinReplica
		}
	}
	if maxReplica == -1 {
		f.logger.Errorf("no such resource spec")
		return fmt.Errorf("no such resource spec")
	}
	realScaleCnt := req.ScaleCnt
	if req.ScaleCnt > 0 {
		if int32(len(currentInsList))+req.ScaleCnt > maxReplica {
			f.logger.Warnf("scale cnt is too large, scale cnt: %d, max replica: %d, current replica: %d", req.ScaleCnt, maxReplica, len(currentInsList))
			realScaleCnt = maxReplica - int32(len(currentInsList))
		} else {
			realScaleCnt = req.ScaleCnt
		}
	} else if req.ScaleCnt < 0 {
		if int32(len(currentInsList))+req.ScaleCnt < minReplica {
			f.logger.Warnf("scale cnt is too small, scale cnt: %d, current replica: %d", req.ScaleCnt, len(currentInsList))
			realScaleCnt = minReplica - int32(len(currentInsList))
		} else {
			realScaleCnt = req.ScaleCnt
		}
	}

	if realScaleCnt > 0 {
		awsServiceNames, err := f.awsClient.BatchCreateInstance(awsTaskDef.TaskFamily, awsTaskDef.TaskRevision, Fargate, realScaleCnt)
		if err != nil {
			f.logger.Errorf("create ecs failed, err: %v", err)
			return err
		}
		var funcInstances []model.FuncInstance
		for _, awsServiceName := range awsServiceNames {
			funcInstances = append(funcInstances, model.FuncInstance{
				AwsServiceName: awsServiceName,
				FunctionName:   req.FunctionName,
				Cpu:            req.Cpu,
				Memory:         req.Memory,
				AwsFamily:      awsTaskDef.TaskFamily,
				AwsRevision:    awsTaskDef.TaskRevision,
				LastStatus:     "NOT_CREATE",
				LaunchType:     int32(Fargate),
			})
		}
		if err := f.mysql.UpdateFuncInstanceBatch(funcInstances); err != nil {
			f.logger.Errorf("UpdateFuncInstanceBatch failed, err: %v", err)
			return err
		}
		f.logger.Infof("scale function %s success, scale cnt: %d", req.FunctionName, realScaleCnt)
		go f.UpdateFunctionStatus(req.FunctionName)
	} else if realScaleCnt < 0 {
		for _, instance := range currentInsList {
			if realScaleCnt >= 0 {
				break
			}
			if err := f.awsClient.DeleteInstance(instance.AwsServiceName); err != nil {
				f.logger.Errorf("delete ecs failed, serviceName: %s, err: %v", instance.AwsServiceName, err)
			}
			if err := f.mysql.DeleteFuncInstanceServiceName(instance.AwsServiceName); err != nil {
				f.logger.Errorf("delete mysql func instance failed, serviceName: %s, err: %v", instance.AwsServiceName, err)
			}
			realScaleCnt++
		}
	}
	return nil
}

func (f *FuncManager) UpdateFunctionStatus(functionName string) {
	startTime := time.Now()

	for {
		funcInstances, err := f.mysql.GetFuncInstanceByFunctionName(functionName)
		if err != nil {
			return
		}
		if isAllReady := f.UpdateFuncInstancesStatus(funcInstances); isAllReady {
			elapsedTime := time.Since(startTime).Seconds()
			f.logger.Infof("%s's all task is ready, cost time: %fs", functionName, elapsedTime)
			return
		}
		time.Sleep(time.Second)
	}

}
func (f *FuncManager) UpdateFuncInstancesStatus(funcInstances []model.FuncInstance) bool {
	// step1: 检查是否所有实例已经就绪
	isAllReady := true
	var taskList []string
	taskMap := make(map[string]int)
	for i, instance := range funcInstances {
		if instance.LastStatus != instance.DesiredStatus {
			isAllReady = false
			tasks, err := f.awsClient.GetAllTasks(instance.AwsServiceName)
			if err != nil {
				f.logger.Errorf("get %s's all tasks failed, err: %v", instance.AwsServiceName, err)
				continue
			}
			if tasks == nil || len(tasks) != 1 {
				funcInstances[i].LastStatus = "NOT_CREATED"
			} else {
				taskList = append(taskList, tasks[0])
				taskMap[tasks[0]] = i
			}
		}
	}
	if isAllReady {
		return isAllReady
	} else if len(taskList) == 0 {
		// task not created
		return false
	}

	// step2: 未就绪的实例获取更新当前状态
	output, err := f.awsClient.DescribeTasks(taskList)
	if err != nil {
		f.logger.Errorf("describe tasks failed, taskList: %v, err: %v", taskList, err)
		return false
	}
	for _, task := range output.Tasks {
		cpu, _ := strconv.Atoi(*task.Cpu)
		memory, _ := strconv.Atoi(*task.Memory)
		publicIpv4, privateIpv4 := "", ""
		if task.Attachments != nil && len(task.Attachments) != 0 {
			for _, d := range task.Attachments[0].Details {
				if d.Name != nil && *d.Name == "networkInterfaceId" {
					publicIpv4, _ = f.awsClient.GetPublicIpv4(*d.Value)
				}
				if d.Name != nil && *d.Name == "privateIPv4Address" {
					privateIpv4 = *d.Value
				}
			}
		}
		instanceIdx := taskMap[*task.TaskArn]
		funcInstances[instanceIdx].AwsTaskArn = *task.TaskArn
		funcInstances[instanceIdx].PrivateIpv4 = privateIpv4
		funcInstances[instanceIdx].PublicIpv4 = publicIpv4
		funcInstances[instanceIdx].Cpu = int32(cpu)
		funcInstances[instanceIdx].Memory = int32(memory)
		funcInstances[instanceIdx].LastStatus = *task.LastStatus
		funcInstances[instanceIdx].DesiredStatus = *task.DesiredStatus
	}
	if err := f.mysql.UpdateFuncInstanceBatch(funcInstances); err != nil {
		f.logger.Errorf("update task status failed, err: %v", err)
	}
	return false
}

func (f *FuncManager) DeleteFunction(req *api.DeleteFunctionRequest) error {
	//step1: check function exist
	_, err := f.mysql.GetFuncMetaDataByFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("get func meta data failed, functionName: %s, err: %v", req.FunctionName, err)
		return err
	}
	funcInstances, err := f.mysql.GetFuncInstanceByFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("get func instance failed, functionName: %s, err: %v", req.FunctionName, err)
		return err
	}

	//step2: delete task
	for _, instance := range funcInstances {
		if err := f.awsClient.DeleteInstance(instance.AwsServiceName); err != nil {
			f.logger.Errorf("delete ecs failed, serviceName: %s, err: %v", instance.AwsServiceName, err)
		}
	}
	err = f.mysql.DeleteFuncInstanceFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("mysql DeleteFuncTaskDataServiceName failed, err:%v", err)
	}
	err = f.mysql.DeleteFuncMetaDataByFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("mysql DeleteFuncTaskDataServiceName failed, err:%v", err)
	}
	return nil
}

func (f *FuncManager) GetAllFunction() (*api.GetAllFunctionsResponse, error) {
	FuncMetaDataList, err := f.mysql.GetFuncMetaDataByCondition(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	resp := &api.GetAllFunctionsResponse{}
	for _, data := range FuncMetaDataList {
		var resSpecList []*api.ResourceSpec
		for _, res := range data.ResSpecList {
			resSpecList = append(resSpecList, &api.ResourceSpec{
				Cpu:        res.Cpu,
				Memory:     res.Memory,
				MinReplica: res.MinReplica,
				MaxReplica: res.MaxReplica,
			})
		}
		resp.Functions = append(resp.Functions, &api.FunctionMetaData{
			FunctionName: data.FunctionName,
			ImageUrl:     data.ImageUrl,
			Resources:    resSpecList,
		})
	}
	return resp, nil
}

func (f *FuncManager) GetFunctionResources(req *api.GetFunctionResourcesRequest) (*api.GetFunctionResourcesResponse, error) {
	taskDataList, err := f.mysql.GetFuncInstanceByFunctionName(req.FunctionName)
	if err != nil {
		return nil, err
	}
	resp := &api.GetFunctionResourcesResponse{
		FunctionName: req.FunctionName,
	}
	for _, taskData := range taskDataList {
		resp.Resources = append(resp.Resources, &api.ResourceStatus{
			PublicIpv4:    taskData.PublicIpv4,
			PrivateIpv4:   taskData.PrivateIpv4,
			Cpu:           taskData.Cpu,
			Memory:        taskData.Memory,
			LaunchType:    taskData.LaunchType,
			LastStatus:    taskData.LastStatus,
			DesiredStatus: taskData.DesiredStatus,
		})
	}
	return resp, nil
}

func (f *FuncManager) FunctionAutoScale() {
	autoScaleStep, autoScaleWindow := config.GetConfig().AutoScaleStep, config.GetConfig().AutoScaleWindow
	windowLen := autoScaleWindow / autoScaleStep
	ticker := time.NewTicker(time.Duration(autoScaleStep) * time.Second)
	defer ticker.Stop()
	hisReqs := list.New()
	for {
		select {
		case <-ticker.C:
			break
		}
		allReqs, err := f.GetAllReq()
		if err != nil {
			f.logger.Errorf("get all req failed, err: %v", err)
			continue
		}
		f.logger.Debugf("current req size: %d", len(allReqs))
		hisReqs.PushFront(allReqs)
		if hisReqs.Len() < windowLen {
			continue
		}
		for hisReqs.Len() > windowLen {
			hisReqs.Remove(hisReqs.Back())
		}
		resDemandMap, err := f.CalResScale(hisReqs)
		if err != nil {
			f.logger.Errorf("cal res demand failed, err: %v", err)
			continue
		}
		for funcName, resDemandList := range resDemandMap {
			for _, resScale := range resDemandList {
				if resScale.ScaleCnt == 0 {
					continue
				}
				err := f.ScaleFunction(&api.ScaleFunctionRequest{
					FunctionName: funcName,
					Cpu:          resScale.Cpu,
					Memory:       resScale.Memory,
					ScaleCnt:     resScale.ScaleCnt,
				})
				if err != nil {
					f.logger.Errorf("scale function failed, resScale: %v, err: %v", resScale, err)
					continue
				}
				f.logger.Infof("scale function success, resScale: %v", resScale)
			}
		}
	}
}

type ResScale struct {
	Cpu      int32
	Memory   int32
	ScaleCnt int32
}

func (f *FuncManager) GetAllReq() ([]*reqscheduler.Request, error) {
	queuedReqs := f.reqQueue.Snapshot()
	scheduledReqs, err := f.mysql.GetReqScheduleInfoByCondition(map[string]interface{}{})
	if err != nil {
		f.logger.Errorf("get mysql req schedule info failed, err: %v", err)
		return nil, err
	}
	var allReq []*reqscheduler.Request
	for _, req := range queuedReqs {
		allReq = append(allReq, &reqscheduler.Request{
			FunctionName:   req.FunctionName,
			RequestID:      req.RequestID,
			RequiredCpu:    req.RequiredCpu,
			RequiredMemory: req.RequiredMemory,
		})
	}
	for _, req := range scheduledReqs {
		allReq = append(allReq, &reqscheduler.Request{
			FunctionName:   req.FunctionName,
			RequestID:      req.ReqId,
			RequiredCpu:    req.RequiredCpu,
			RequiredMemory: req.RequiredMemory,
		})
	}
	return allReq, nil
}

func (f *FuncManager) CalResScale(hisReqQueue *list.List) (map[string][]*ResScale, error) {
	funcMetaDataList, err := f.mysql.GetFuncMetaDataByCondition(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	funcMetaDataMap := make(map[string]model.FuncMetaData)
	// [ts]map[funcName][cpu+memory]required_cnt
	resDemandMap := make([]map[string]map[[2]int32]float64, hisReqQueue.Len())
	currentResMap := make(map[string]map[[2]int32]int32)
	for i := 0; i < hisReqQueue.Len(); i++ {
		resDemandMap[i] = make(map[string]map[[2]int32]float64)
	}
	for _, metaData := range funcMetaDataList {
		sort.Slice(metaData.ResSpecList, func(i, j int) bool {
			if metaData.ResSpecList[i].Cpu != metaData.ResSpecList[j].Cpu {
				return metaData.ResSpecList[i].Cpu < metaData.ResSpecList[j].Cpu
			} else {
				return metaData.ResSpecList[i].Memory < metaData.ResSpecList[j].Memory
			}
		})
		funcMetaDataMap[metaData.FunctionName] = metaData
		for i := 0; i < hisReqQueue.Len(); i++ {
			resDemandMap[i][metaData.FunctionName] = make(map[[2]int32]float64)
		}
		currentResMap[metaData.FunctionName] = make(map[[2]int32]int32)
	}

	allInsList, err := f.mysql.GetFuncInstanceByCondition(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	for _, ins := range allInsList {
		currentResMap[ins.FunctionName][[2]int32{ins.Cpu, ins.Memory}]++
	}

	it := hisReqQueue.Front()
	for i := 0; i < hisReqQueue.Len(); i++ {
		reqQueue := it.Value.([]*reqscheduler.Request)
		for _, req := range reqQueue {
			if metaData, ok := funcMetaDataMap[req.FunctionName]; ok {
				maxRes := metaData.ResSpecList[len(metaData.ResSpecList)-1]
				resCpu, resMemory, requireCnt := maxRes.Cpu, maxRes.Memory, max(float64(req.RequiredCpu)/float64(maxRes.Cpu), float64(req.RequiredMemory)/float64(maxRes.Memory))
				for _, res := range metaData.ResSpecList {
					if req.RequiredCpu <= res.Cpu && req.RequiredMemory <= res.Memory {
						resCpu, resMemory, requireCnt = res.Cpu, res.Memory, max(float64(req.RequiredCpu)/float64(res.Cpu), float64(req.RequiredMemory)/float64(res.Memory))
						break
					}
				}
				resDemandMap[i][metaData.FunctionName][[2]int32{resCpu, resMemory}] += requireCnt
			}
		}
		it = it.Next()
	}

	avgResDemandMap := make(map[string]map[[2]int32]float64)
	for _, metaData := range funcMetaDataList {
		avgResDemandMap[metaData.FunctionName] = make(map[[2]int32]float64)
		for i := 0; i < hisReqQueue.Len(); i++ {
			for res, demand := range resDemandMap[i][metaData.FunctionName] {
				avgResDemandMap[metaData.FunctionName][res] += demand / float64(hisReqQueue.Len())
			}
		}
	}

	ret := make(map[string][]*ResScale)
	for funcName, metaData := range funcMetaDataMap {
		ret[funcName] = []*ResScale{}
		for _, res := range metaData.ResSpecList {
			ret[funcName] = append(ret[funcName], &ResScale{
				Cpu:      res.Cpu,
				Memory:   res.Memory,
				ScaleCnt: min(res.MaxReplica, max(res.MinReplica, int32(math.Round(avgResDemandMap[funcName][[2]int32{res.Cpu, res.Memory}])))) - currentResMap[funcName][[2]int32{res.Cpu, res.Memory}],
			})
		}
	}
	//for funcName, resMap := range currentResMap {
	//	ret[funcName] = []*ResScale{}
	//	for res, cnt := range resMap {
	//		ret[funcName] = append(ret[funcName], &ResScale{
	//			Cpu:      res[0],
	//			Memory:   res[1],
	//			ScaleCnt: int32(math.Round(avgResDemandMap[funcName][res])) - cnt,
	//		})
	//	}
	//}
	f.logger.Debugf("res demand: %v, res scale: %s", avgResDemandMap, utils.GetJson(ret))
	return ret, nil
}
