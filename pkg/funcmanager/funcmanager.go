/*
@author: panfengguo
@since: 2024/11/6
@decs: decs
*/
package funcmanager

import (
	"fmt"
	"github.com/AgentGuo/faas/api"
	"github.com/AgentGuo/faas/pkg/logger"
	"github.com/AgentGuo/faas/pkg/storage"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type FuncManager struct {
	awsClient   *AwsClient
	mysqlClient *storage.MysqlClient
	logger      *logrus.Logger
}

var (
	funcManager *FuncManager
	once        sync.Once
)

func NewFuncManager() *FuncManager {
	once.Do(func() {
		funcManager = &FuncManager{
			awsClient:   NewAwsClient(),
			mysqlClient: storage.NewMysqlClient(),
			logger:      logger.GetLogger(),
		}
	})
	return funcManager
}

func (f *FuncManager) CreateFunction(req *api.CreateFunctionRequest) error {
	// step1: check if function has been created
	if hasCreated, err := funcManager.mysqlClient.HasFuncMetaDataByFunctionName(req.FunctionName); err != nil || hasCreated {
		if hasCreated {
			return fmt.Errorf("function has been created")
		} else {
			return err
		}
	}

	// step2: create task definition
	var resourceSpecList []storage.ResourceSpec
	for _, res := range req.Resources {
		family, revision, err := f.awsClient.RegTaskDef(req.FunctionName, res.Cpu, res.Memory, req.ImageUrl)
		if err != nil {
			return err
		}
		resourceSpecList = append(resourceSpecList, storage.ResourceSpec{
			Cpu:               res.Cpu,
			Memory:            res.Memory,
			MinReplica:        res.MinReplica,
			MaxReplica:        res.MaxReplica,
			EnableAutoScaling: res.EnableAutoScaling,
			ServiceName:       fmt.Sprintf("%s_v%d", family, revision),
			Family:            family,
			Revision:          revision,
		})
	}
	err := funcManager.mysqlClient.CreateFuncMetaData(&storage.FuncMetaData{
		FunctionName: req.FunctionName,
		ImageUrl:     req.ImageUrl,
		ResSpecList:  resourceSpecList,
	})
	if err != nil {
		f.logger.Errorf("create func meta data failed, err: %v", err)
		return err
	}

	// step3: create task
	for _, res := range resourceSpecList {
		if _, err := f.awsClient.CreateECS(res.Family, res.Revision, res.MinReplica); err != nil {
			f.logger.Errorf("create ecs failed, err: %v", err)
			return err
		}
	}
	go f.UpdateTaskStatusRoutine(req.FunctionName)
	return nil
}

func (f *FuncManager) UpdateTaskStatusRoutine(functionName string) {
	startTime := time.Now()
	for {
		funcMeta, err := f.mysqlClient.GetFuncMetaDataByFunctionName(functionName)
		if err != nil {
			return
		}
		if isAllReady, err := f.UpdateTaskStatus(funcMeta); err != nil {
			f.logger.Errorf("update task status failed, err: %v", err)
		} else if isAllReady {
			elapsedTime := time.Since(startTime).Seconds()
			f.logger.Infof("%s's all task is ready, cost time: %fs", functionName, elapsedTime)
			return
		}
		time.Sleep(time.Second)
	}
}

func (f *FuncManager) UpdateTaskStatus(metaData *storage.FuncMetaData) (bool, error) {
	isAllReady := true
	for _, resSpec := range metaData.ResSpecList {
		// step1: get tasks status
		tasks, err := f.awsClient.GetAllTasks(resSpec.ServiceName)
		if err != nil {
			return false, err
		}
		output, err := f.awsClient.DescribeTasks(tasks)
		if err != nil {
			return false, err
		}

		// step2: update tasks status
		var updateTask []storage.FuncTaskData
		if output == nil || output.Tasks == nil {
			if resSpec.MinReplica != 0 {
				isAllReady = false
			}
			continue
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
			updateTask = append(updateTask, storage.FuncTaskData{
				TaskArn:       *task.TaskArn,
				ServiceName:   resSpec.ServiceName,
				FunctionName:  metaData.FunctionName,
				PrivateIpv4:   privateIpv4,
				PublicIpv4:    publicIpv4,
				Cpu:           int32(cpu),
				Memory:        int32(memory),
				Family:        resSpec.Family,
				Revision:      resSpec.Revision,
				LastStatus:    *task.LastStatus,
				DesiredStatus: *task.DesiredStatus,
				LaunchType:    string(task.LaunchType),
			})
			if *task.LastStatus != *task.DesiredStatus {
				isAllReady = false
			}
		}
		if int32(len(output.Tasks)) < resSpec.MinReplica {
			isAllReady = false
		}
		currentTask, err := f.mysqlClient.GetFuncTaskDataByServiceName(resSpec.ServiceName)
		if err != nil {
			return false, err
		}
		if len(currentTask) > len(updateTask) {
			_ = f.mysqlClient.DeleteFuncTaskDataServiceName(resSpec.ServiceName)
		}
		if err := f.mysqlClient.UpdateFuncTaskDataBatch(updateTask); err != nil {
			return false, err
		}
	}
	return isAllReady, nil
}

func (f *FuncManager) DeleteFunction(req *api.DeleteFunctionRequest) error {
	//step1: get funcMetaData
	funcMetaData, err := funcManager.mysqlClient.GetFuncMetaDataByFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("get func meta data failed, functionName: %s, err: %v", req.FunctionName, err)
		return err
	}

	//step2: delete task
	for _, resSpec := range funcMetaData.ResSpecList {
		if err := f.awsClient.DeleteECS(resSpec.ServiceName); err != nil {
			f.logger.Errorf("delete ecs failed, serviceName: %s, err: %v", resSpec.ServiceName, err)
		}
	}
	err = f.mysqlClient.DeleteFuncTaskDataFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("mysql DeleteFuncTaskDataServiceName failed, err:%v", err)
	}
	err = f.mysqlClient.DeleteFuncMetaDataByFunctionName(req.FunctionName)
	if err != nil {
		f.logger.Errorf("mysql DeleteFuncTaskDataServiceName failed, err:%v", err)
	}
	return nil
}

func (f *FuncManager) GetAllFunction() (*api.GetAllFunctionsResponse, error) {
	FuncMetaDataList, err := f.mysqlClient.GetFuncMetaDataByCondition(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	resp := &api.GetAllFunctionsResponse{}
	for _, data := range FuncMetaDataList {
		var resSpecList []*api.ResourceSpec
		for _, res := range data.ResSpecList {
			resSpecList = append(resSpecList, &api.ResourceSpec{
				Cpu:               res.Cpu,
				Memory:            res.Memory,
				MinReplica:        res.MinReplica,
				MaxReplica:        res.MaxReplica,
				EnableAutoScaling: res.EnableAutoScaling,
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
	taskDataList, err := f.mysqlClient.GetFuncTaskDataByFunctionName(req.FunctionName)
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
