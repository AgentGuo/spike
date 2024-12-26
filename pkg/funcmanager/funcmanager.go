/*
@author: panfengguo
@since: 2024/11/6
@decs: decs
*/
package funcmanager

import (
	"fmt"
	"github.com/AgentGuo/spike/api"
	"github.com/AgentGuo/spike/pkg/logger"
	"github.com/AgentGuo/spike/pkg/storage"
	"github.com/AgentGuo/spike/pkg/storage/model"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type FuncManager struct {
	awsClient *AwsClient
	mysql     *storage.Mysql
	logger    *logrus.Logger
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
		}
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
