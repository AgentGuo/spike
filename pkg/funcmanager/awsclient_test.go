/*
@author: panfengguo
@since: 2024/11/16
@desc: desc
*/
package funcmanager

import (
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
		{"test1", args{"faas_test", 1024, 3072, "013072238852.dkr.ecr.cn-north-1.amazonaws.com.cn/agentguo/spike-java-worker:1.0"}, false},
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

func TestAwsClient_CreateESC(t *testing.T) {
	type args struct {
		familyName string
		revision   int32
		replicas   int32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"faas_test", 1, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAwsClient()
			if _, err := a.CreateECS(tt.args.familyName, tt.args.revision, tt.args.replicas); (err != nil) != tt.wantErr {
				t.Errorf("CreateECS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAwsClient_UpdateECSReplicas(t *testing.T) {
	type args struct {
		serviceName string
		replicas    int32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"faas_test_v1", 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAwsClient()
			if err := a.UpdateECSReplicas(tt.args.serviceName, tt.args.replicas); (err != nil) != tt.wantErr {
				t.Errorf("UpdateECSReplicas() error = %v, wantErr %v", err, tt.wantErr)
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
		{"test1", args{"test_v14"}, false},
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
		{"test1", args{"eni-0580f1a35d9335b77"}, false},
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
