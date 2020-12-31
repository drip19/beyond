package main

import (
	"fmt"
	"github.com/ymer/mydemo/config"
	"github.com/ymer/mydemo/pkg/log"
	"github.com/ymer/mydemo/pkg/util"
	"github.com/ymer/mydemo/rpc"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/jessevdk/go-flags"
)

var (
	version = "dev"
	date    = ""
	commit  = ""
)

var cfg = &config.Config{}

func main() {
	fmt.Println("start...")

	if _, err := flag.ParseArgs(cfg, os.Args); err != nil {
		code := 1
		if fe, ok := err.(*flag.Error); ok {
			if fe.Type == flag.ErrHelp {
				code = 0
			} else {
				log.Root.Error(err)
			}
		}
		os.Exit(code)
	}

	if err := cfg.Load(); err != nil {
		log.Root.Error(err)
	}

	if err := cfg.Verify(); err != nil {
		fmt.Println(util.ToIndentString(cfg))
		log.Root.Fatal(err)
	}
	if err := cfg.Save(); err != nil {
		log.Root.Error(err)
	}

	if cfg.Verbose {
		cfg.LogLevel = "debug"
	}
	_ = log.Setup(cfg.LogDir(), cfg.LogLevel)

	logger := log.NewLogger("main")
	logger.Info(util.ToIndentString(cfg))

	server := rpc.NewServer(cfg)
	if err := server.Start(); err != nil {
		logger.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

}
