package mfunc

import (
	qlcchain "github.com/qlcchain/qlc-go-sdk"
	"github.com/ymer/mydemo/config"
	"github.com/ymer/mydemo/pkg/log"
	"time"
)

func Tfunc(cfg *config.Config) {
	logger := log.NewLogger("tfunc")
	logger.Info("endpoint: ", cfg.Endpoint)
	client, err := qlcchain.NewQLCClient(cfg.Endpoint)
	if err != nil || client == nil {
		logger.Fatal(err)
	} else {
		defer client.Close()
	}

	for {
		time.Sleep(10 * time.Second)
		p, err := client.Ledger.Tokens()
		logger.Info(time.Now().Second(), p, err)
	}
}
