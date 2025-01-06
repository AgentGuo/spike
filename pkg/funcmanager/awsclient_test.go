/*
@author: panfengguo
@since: 2024/11/16
@desc: desc
*/
package funcmanager

import (
	"github.com/AgentGuo/spike/pkg/constants"
	"testing"
)

func TestAwsClient_RegTaskDef(t *testing.T) {
	type args struct {
		functionName string
		cpu          int32
		memory       int32
		imageUrl     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"spike_test", 1024, 3072, "013072238852.dkr.ecr.cn-north-1.amazonaws.com.cn/agentguo/spike-java-worker:1.0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAwsClient()
			if _, _, err := a.RegTaskDef(tt.args.functionName, tt.args.cpu, tt.args.memory, tt.args.imageUrl); (err != nil) != tt.wantErr {
				t.Errorf("RegTaskDef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAwsClient_GetAllTasks(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test_fargate", args{"pixels-worker-spike_546127393488676365"}, false},
		{"test_fargate_spot", args{"pixels-worker-spike_546127394595972621"}, false},
		{"test_ec2", args{"pixels-worker-spike_546127776462185997"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAwsClient()
			got, err := a.GetAllTasks(tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
			_, err = a.DescribeTasks(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DescribeTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAwsClient_DescribeDescribeNetworkInterfaces(t *testing.T) {
	type args struct {
		interfaceId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test_fargate", args{"eni-0aa95680207aef0b9"}, false},
		{"test_fargate_spot", args{"eni-0899dd3c4119a3791"}, false},
		{"test_ec2", args{"eni-058b42ef82434c04e"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAwsClient()
			if got, err := a.GetPublicIpv4(tt.args.interfaceId); (err != nil) != tt.wantErr {
				t.Errorf("DescribeDescribeNetworkInterfaces() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestAwsClient_CreateInstance(t *testing.T) {
	type args struct {
		familyName   string
		revision     int32
		instanceType constants.InstanceType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"pixels-worker-spike", 6, constants.EC2}, false},
		{"test2", args{"pixels-worker-spike", 6, constants.Fargate}, false},
		{"test3", args{"pixels-worker-spike", 6, constants.FargateSpot}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAwsClient()
			got, err := a.CreateInstance(tt.args.familyName, tt.args.revision, tt.args.instanceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("CreateInstance() got = %v", got)
		})
	}
}
