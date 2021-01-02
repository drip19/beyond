package jsonrpc

import (
	"github.com/drip/beyond/config"
	"github.com/drip/beyond/pkg/log"
	"go.uber.org/zap"
)

type RPCService struct {
	rpc    *RPC
	logger *zap.SugaredLogger
}

func NewRPCService(cfg *config.Config) (*RPCService, error) {
	logger := log.NewLogger("rpc_service")
	rpc, err := NewRPC(cfg)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info("grpc started")
	return &RPCService{rpc: rpc, logger: logger}, nil
}

func (r *RPCService) Start() error {
	return r.rpc.StartRPC()
}

func (r *RPCService) Stop() {
	r.rpc.StopRPC()
	r.logger.Info("wrapper grpc stopped")
}
