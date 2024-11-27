/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package funcmanager

import (
	"context"
	"fmt"
	"github.com/AgentGuo/faas/cmd/server/config"
	"github.com/AgentGuo/faas/pkg/logger"
	"github.com/AgentGuo/faas/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/sirupsen/logrus"
)

type AwsClient struct {
	awsCfg         *aws.Config
	ecsClient      *ecs.Client
	ec2Client      *ec2.Client
	cluster        string
	subnets        []string
	securityGroups []string
	logger         *logrus.Logger
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
		cluster:        config.GetConfig().AwsFassCluster,
		subnets:        config.GetConfig().AwsSubnets,
		securityGroups: config.GetConfig().AwsSecurityGroups,
		logger:         logger.GetLogger(),
	}
}

func (a *AwsClient) CreateECS(familyName string, revision int32, replicas int32) (string, error) {
	serviceName := fmt.Sprintf("%s_v%d", familyName, revision)
	output, err := a.ecsClient.CreateService(context.TODO(), &ecs.CreateServiceInput{
		ServiceName:    aws.String(serviceName),
		Cluster:        aws.String(a.cluster),
		DesiredCount:   aws.Int32(replicas),
		TaskDefinition: aws.String(fmt.Sprintf("%s:%d", familyName, revision)),
		//LaunchType:     types.LaunchTypeFargate,
		//LaunchType: types.LaunchTypeEc2,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        a.subnets,
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: a.securityGroups,
			},
		},
	})
	if err != nil {
		a.logger.Errorf("failed to create ECS, err: %v, resp: %s", err, utils.GetJson(output))
		return "", err
	}
	a.logger.Debugf("create ECS success, resp: %s", utils.GetJson(output))
	return serviceName, nil
}

func (a *AwsClient) UpdateECSReplicas(serviceName string, replicas int32) error {
	output, err := a.ecsClient.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
		Service:      aws.String(serviceName),
		Cluster:      aws.String(a.cluster),
		DesiredCount: aws.Int32(replicas),
	})
	if err != nil {
		a.logger.Errorf("failed to update ECS replicas, err: %v, resp: %#v", err, output)
	}
	a.logger.Debugf("update ECS replicas success, resp: %#v", output)
	return nil
}

func (a *AwsClient) GetAllTasks(serviceName string) ([]string, error) {
	listTaskOutput, err := a.ecsClient.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster:     aws.String(a.cluster),
		ServiceName: aws.String(serviceName),
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

func (a *AwsClient) DeleteECS(serviceName string) error {
	output, err := a.ecsClient.DeleteService(context.TODO(), &ecs.DeleteServiceInput{
		Service: aws.String(serviceName),
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

func (a *AwsClient) RegTaskDef(functionName string, cpu int32, memory int32, imageUrl string) (string, int32, error) {
	output, err := a.ecsClient.RegisterTaskDefinition(context.TODO(), &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []types.ContainerDefinition{
			{
				Image:  aws.String(imageUrl),
				Cpu:    cpu,
				Memory: aws.Int32(memory),
				Name:   aws.String("faas_worker"),
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
						"awslogs-group":         "pixels-worker-faas",
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
