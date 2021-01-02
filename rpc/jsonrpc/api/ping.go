package api

import (
	"github.com/drip/beyond/config"
	qlcchain "github.com/qlcchain/qlc-go-sdk"
)

type PingApi struct {
	cfg    *config.Config
	client *qlcchain.QLCClient
}

func NewPingApi(cfg *config.Config) *PingApi {
	client, err := qlcchain.NewQLCClient(cfg.Endpoint)
	if err != nil || client == nil {
		panic(err)
	}
	return &PingApi{
		client: client,
		cfg:    cfg,
	}
}

func (p *PingApi) Info() (string, error) {
	return "ping.info", nil
}

func (p *PingApi) State() (bool, error) {
	_, err := p.client.Ledger.Tokens()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
