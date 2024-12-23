/*
@author: panfengguo
@since: 2024/11/6
@desc: desc
*/
package apiserver

import (
	"context"
	"fmt"
	"github.com/AgentGuo/spike/api"
	"github.com/AgentGuo/spike/cmd/server/config"
	"github.com/AgentGuo/spike/pkg/funcdispatch"
	"github.com/AgentGuo/spike/pkg/funcmanager"
	"github.com/AgentGuo/spike/pkg/logger"
	"github.com/sirupsen/logrus"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	logger *logrus.Logger
	api.UnimplementedSpikeServiceServer
	funcManager  *funcmanager.FuncManager
	funcDispatch *funcdispatch.FuncDispatch
}

func (s *server) CallFunction(ctx context.Context, req *api.CallFunctionRequest) (*api.CallFunctionResponse, error) {
	resp, err := s.funcDispatch.CallFunction(req)
	if err != nil {
		s.logger.Errorf("call function %s failed, err: %v", req.FunctionName, err)
	}
	return resp, err
}

func (s *server) CreateFunction(ctx context.Context, req *api.CreateFunctionRequest) (*api.CreateFunctionResponse, error) {
	err := s.funcManager.CreateFunction(req)
	if err != nil {
		return nil, err
	}
	return &api.CreateFunctionResponse{Code: 0, Message: "Function added"}, nil
}

func (s *server) DeleteFunction(ctx context.Context, req *api.DeleteFunctionRequest) (*api.DeleteFunctionResponse, error) {
	err := s.funcManager.DeleteFunction(req)
	if err != nil {
		return nil, err
	}
	return &api.DeleteFunctionResponse{Code: 0, Message: "Function deleted"}, nil
}

func (s *server) GetAllFunctions(context.Context, *api.Empty) (*api.GetAllFunctionsResponse, error) {
	resp, err := s.funcManager.GetAllFunction()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *server) GetFunctionResources(ctx context.Context, req *api.GetFunctionResourcesRequest) (*api.GetFunctionResourcesResponse, error) {
	return s.funcManager.GetFunctionResources(req)
}

func StartApiServer() {
	address := fmt.Sprintf("%s:%d", config.GetConfig().ServerIp, config.GetConfig().ServerPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterSpikeServiceServer(grpcServer, &server{
		logger:       logger.GetLogger(),
		funcManager:  funcmanager.NewFuncManager(),
		funcDispatch: funcdispatch.NewFuncDispatch(),
	})

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Printf("gRPC server is running on %s\n", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
