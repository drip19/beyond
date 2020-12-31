package apis

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	qlcchain "github.com/qlcchain/qlc-go-sdk"
	"github.com/ymer/mydemo/config"
	"github.com/ymer/mydemo/pkg/log"
	pb "github.com/ymer/mydemo/rpc/proto"
	"go.uber.org/zap"
)

type PingApi struct {
	cfg    *config.Config
	client *qlcchain.QLCClient
	logger *zap.SugaredLogger
}

func NewPingApi(cfg *config.Config) *PingApi {
	client, err := qlcchain.NewQLCClient(cfg.Endpoint)
	if err != nil || client == nil {
		panic(err)
	}
	return &PingApi{
		cfg:    cfg,
		client: client,
		logger: log.NewLogger("api/ping"),
	}
}

func (p *PingApi) Info(ctx context.Context, empty *empty.Empty) (*pb.String, error) {
	return &pb.String{
		Value: "ping.info",
	}, nil
}

func (p *PingApi) Status(ctx context.Context, e *empty.Empty) (*pb.Boolean, error) {
	_, err := p.client.Ledger.Tokens()
	if err != nil {
		return &pb.Boolean{Value: false}, err
	} else {
		return &pb.Boolean{Value: true}, nil
	}
}
