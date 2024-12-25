/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package funcmanager

import (
	"context"
	"fmt"
	"github.com/AgentGuo/spike/cmd/server/config"
	"github.com/AgentGuo/spike/pkg/logger"
	"github.com/AgentGuo/spike/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/sirupsen/logrus"
	"sync"
)

type InstanceType int

const (
	EC2 InstanceType = iota
	Fargate
	FargateSpot
)

type AwsClient struct {
	awsCfg            *aws.Config
	ecsClient         *ecs.Client
	ec2Client         *ec2.Client
	cluster           string
	subnets           []string
	securityGroups    []string
	logger            *logrus.Logger
	CapacityProviders map[InstanceType]string
}

func NewAwsClient() *AwsClient {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	return &AwsClient{
		awsCfg:         &cfg,
		ecsClient:      ecs.NewFromConfig(cfg),
		ec2Client:      ec2.NewFromConfig(cfg),
		cluster:        config.GetConfig().AwsCluster,
		subnets:        config.GetConfig().AwsSubnets,
		securityGroups: config.GetConfig().AwsSecurityGroups,
		logger:         logger.GetLogger(),
		CapacityProviders: map[InstanceType]string{
			EC2:         config.GetConfig().EC2Provider,
			Fargate:     "FARGATE",
			FargateSpot: "FARGATE_SPOT",
		},
	}
}

func (a *AwsClient) GenServiceName(awsFamilyName string) string {
	return fmt.Sprintf("%s_%d", awsFamilyName, utils.GetSonyFlakeInstance().GenerateID())
}

func (a *AwsClient) BatchCreateInstance(awsFamilyName string, awsRevision int32, instanceType InstanceType, replicas int32) ([]string, error) {
	var wg sync.WaitGroup
	awsServiceNames := make([]string, replicas)
	errors := make([]error, replicas)

	for i := int32(0); i < replicas; i++ {
		wg.Add(1)
		go func(index int32) {
			defer wg.Done()
			awsServiceName, err := a.CreateInstance(awsFamilyName, awsRevision, instanceType)
			awsServiceNames[index] = awsServiceName
			errors[index] = err
		}(i)
	}

	wg.Wait()
	var retErr error
	for _, err := range errors {
		if err != nil {
			retErr = fmt.Errorf("failed to create some instances: %v", errors)
		}
	}
	if retErr != nil {
		for i, err := range errors {
			if err != nil {
				if delErr := a.DeleteInstance(awsServiceNames[i]); delErr != nil {
					a.logger.Errorf("failed to delete instance %s, err: %v", awsServiceNames[i], delErr)
				}
			}
		}
		return nil, retErr
	}
	return awsServiceNames, nil
}

// CreateInstance 为了方便管理，避免出现热实例缩容的情况，所以创建一个一个service，
// 而不是一个service下创建多个replicas，这样可以准确控制扩缩容的实例
func (a *AwsClient) CreateInstance(awsFamilyName string, awsRevision int32, instanceType InstanceType) (string, error) {
	awsServiceName := a.GenServiceName(awsFamilyName)
	assignPublicIp := types.AssignPublicIpEnabled
	if instanceType == EC2 {
		assignPublicIp = types.AssignPublicIpDisabled
	}
	output, err := a.ecsClient.CreateService(context.TODO(), &ecs.CreateServiceInput{
		ServiceName: aws.String(awsServiceName),
		CapacityProviderStrategy: []types.CapacityProviderStrategyItem{
			{
				CapacityProvider: aws.String(a.CapacityProviders[instanceType]),
				Base:             0,
				Weight:           1,
			},
		},
		Cluster:        aws.String(a.cluster),
		DesiredCount:   aws.Int32(1),
		TaskDefinition: aws.String(fmt.Sprintf("%s:%d", awsFamilyName, awsRevision)),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        a.subnets,
				AssignPublicIp: assignPublicIp,
				SecurityGroups: a.securityGroups,
			},
		},
	})
	if err != nil {
		a.logger.Errorf("failed to create instance, err: %v, resp: %s", err, utils.GetJson(output))
		return "", err
	}
	a.logger.Debugf("create instance success, resp: %s", utils.GetJson(output))
	return awsServiceName, nil
}

func (a *AwsClient) DeleteInstance(awsServiceName string) error {
	output, err := a.ecsClient.DeleteService(context.TODO(), &ecs.DeleteServiceInput{
		Service: aws.String(awsServiceName),
		Cluster: aws.String(a.cluster),
		Force:   aws.Bool(true),
	})
	if err != nil {
		a.logger.Errorf("failed to delete ECS, err: %v, resp: %#v", err, output)
		return err
	}
	a.logger.Debugf("delete ECS success, resp: %#v", output)
	return nil
}

func (a *AwsClient) GetAllTasks(awsServiceName string) ([]string, error) {
	listTaskOutput, err := a.ecsClient.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster:     aws.String(a.cluster),
		ServiceName: aws.String(awsServiceName),
	})
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}
	return listTaskOutput.TaskArns, nil
}

func (a *AwsClient) DescribeTasks(tasks []string) (*ecs.DescribeTasksOutput, error) {
	if len(tasks) == 0 {
		return nil, nil
	}
	output, err := a.ecsClient.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Tasks:   tasks,
		Cluster: aws.String(a.cluster),
	})
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}
	a.logger.Debug("output:", utils.GetJson(output))
	return output, nil
}

func (a *AwsClient) RegTaskDef(functionName string, cpu int32, memory int32, imageUrl string) (string, int32, error) {
	output, err := a.ecsClient.RegisterTaskDefinition(context.TODO(), &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []types.ContainerDefinition{
			{
				Image:  aws.String(imageUrl),
				Cpu:    cpu,
				Memory: aws.Int32(memory),
				Name:   aws.String("spike_worker"),
				PortMappings: []types.PortMapping{{
					AppProtocol:   types.ApplicationProtocolGrpc,
					ContainerPort: aws.Int32(50052),
					HostPort:      aws.Int32(50052),
					Name:          aws.String("invoke_port"),
					Protocol:      types.TransportProtocolTcp,
				}},
				LogConfiguration: &types.LogConfiguration{
					LogDriver: types.LogDriverAwslogs,
					Options: map[string]string{
						"awslogs-region":        "cn-north-1",
						"awslogs-group":         "pixels-worker-spike",
						"awslogs-stream-prefix": "ecs",
					},
					SecretOptions: nil,
				},
			},
		},
		Family:      aws.String(functionName),
		Cpu:         aws.String(fmt.Sprintf("%d", cpu)),
		Memory:      aws.String(fmt.Sprintf("%d", memory)),
		NetworkMode: types.NetworkModeAwsvpc,
		RequiresCompatibilities: []types.Compatibility{
			types.CompatibilityEc2,
			types.CompatibilityFargate,
		},
		RuntimePlatform: &types.RuntimePlatform{
			CpuArchitecture:       types.CPUArchitectureX8664,
			OperatingSystemFamily: types.OSFamilyLinux,
		},
		ExecutionRoleArn: aws.String(config.GetConfig().TaskRole),
		TaskRoleArn:      aws.String(config.GetConfig().TaskRole),
	})
	if err != nil {
		a.logger.Errorf("failed to register task definition, err: %v, resp: %s", err, utils.GetJson(output))
		return "", 0, err
	}
	a.logger.Debugf("register task definition success, resp: %s", utils.GetJson(output))
	return *output.TaskDefinition.Family, output.TaskDefinition.Revision, nil
}

func (a *AwsClient) GetPublicIpv4(interfaceIds string) (string, error) {
	output, err := a.ec2Client.DescribeNetworkInterfaces(context.TODO(), &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []string{interfaceIds},
	})
	if err != nil {
		a.logger.Errorf("failed to DescribeDescribeNetworkInterfaces, err: %v, resp: %s", err, utils.GetJson(output))
		return "", err
	}
	a.logger.Debugf("DescribeDescribeNetworkInterfaces success, resp: %s", utils.GetJson(output))
	publicIpv4 := ""
	for _, item := range output.NetworkInterfaces {
		if item.Association != nil && item.Association.PublicIp != nil {
			publicIpv4 = *item.Association.PublicIp
		}
		break
	}
	return publicIpv4, nil
}
