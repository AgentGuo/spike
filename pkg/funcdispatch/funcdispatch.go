/*
@author: panfengguo
@since: 2024/11/17
@desc: desc
*/
package funcdispatch

import (
	"context"
	"fmt"
	"github.com/AgentGuo/faas/api"
	"github.com/AgentGuo/faas/pkg/funcdispatch/proto"
	"github.com/AgentGuo/faas/pkg/logger"
	"github.com/AgentGuo/faas/pkg/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"time"
)

// Request 表示一个请求
type Request struct {
	ID      int64
	Payload string
}

type FuncDispatch struct {
	mysqlClient *storage.MysqlClient
	logger      *logrus.Logger
}

func NewFuncDispatch() *FuncDispatch {
	mysqlClient := storage.NewMysqlClient()
	return &FuncDispatch{
		mysqlClient: mysqlClient,
		logger:      logger.GetLogger(),
	}
}

func (f *FuncDispatch) AddRequest(req *api.CallFunctionRequest) {
	// TODO
}

func (f *FuncDispatch) CallFunction(req *api.CallFunctionRequest) (*api.CallFunctionResponse, error) {
	taskData, err := f.mysqlClient.GetFuncTaskDataByCondition(map[string]interface{}{
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
	funcServiceClient := proto.NewFunctionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	funcServiceResp, err := funcServiceClient.CallFunction(ctx, &proto.FunctionRequest{
		Payload:   req.Payload,
		RequestId: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return &api.CallFunctionResponse{ErrorCode: funcServiceResp.ErrorCode, Payload: funcServiceResp.Payload}, nil
}
