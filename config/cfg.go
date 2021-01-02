package config

import (
	"encoding/json"
	"github.com/drip/beyond/pkg/util"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/validator.v2"
)

const (
	cfgDir    = "gbeyond"
	nixCfgDir = ".gbeyond"
)

type Config struct {
	Verbose  bool     `json:"verbose" long:"verbose" description:"show verbose debug information"`
	LogLevel string   `json:"logLevel" long:"level" description:"log level" default:"debug"` //info,warn,debug.
	Names    []string `json:"names"  validate:"min=0"`
	Endpoint string   `json:"endpoint" long:"endpoint" description:"endpoint" default:"ws://127.0.0.1:29736"`
	GRPCCfg  *GRPCCfg `json:"grpc" validate:"nonnil"`
	RPCCfg   *RPCCfg  `json:"rpc" validate:"nonnil"`
}

type GRPCCfg struct {
	// TCP or UNIX socket address for the RPC server to listen on
	ListenAddress string `json:"listenAddress" long:"listenAddress" description:"RPC server listen address" default:"tcp://0.0.0.0:29705"`
	// TCP or UNIX socket address for the gRPC server to listen on
	GRPCListenAddress  string   `json:"gRPCListenAddress" long:"grpcAddress" description:"GRPC server listen address" default:"tcp://0.0.0.0:29706"`
	CORSAllowedOrigins []string `json:"allowedOrigins" long:"allowedOrigins" description:"AllowedOrigins of CORS" default:"*"`
}

type RPCCfg struct {
	Enable           bool     `json:"rpcEnabled"`
	HTTPEndpoint     string   `json:"httpEndpoint" long:"httpEndpoint" default:"tcp://0.0.0.0:29707"`
	HTTPCors         []string `json:"httpCors" default:"*"`
	HttpVirtualHosts []string `json:"httpVirtualHosts" default:"*"`
	HTTPEnabled      bool     `json:"httpEnabled" `

	WSEndpoint string `json:"webSocketEndpoint" long:"WSEndpoint" default:"tcp://0.0.0.0:29708"`
	WSEnabled  bool   `json:"webSocketEnabled" `

	IPCEndpoint string `json:"ipcEndpoint"`
	IPCEnabled  bool   `json:"ipcEnabled" `
}

func (c *Config) LogDir() string {
	return filepath.Join(DefaultDataDir(), "log", time.Now().Format("2006-01-02T15-04"))
}

func (c *Config) Database() string {
	dir := filepath.Join(DefaultDataDir(), "db")
	_ = util.CreateDirIfNotExist(dir)

	return filepath.Join(dir, "actions.db")
}

func (c *Config) Load() error {
	f := filepath.Join(DefaultDataDir(), "config.json")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		if data, err := ioutil.ReadFile(f); err == nil {
			cfg := &Config{}
			if err := json.Unmarshal(data, cfg); err == nil {

			} else {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (c *Config) Save() error {
	f := filepath.Join(DefaultDataDir(), "config.json")
	s := util.ToIndentString(c)
	//data, _ := json.Marshal(c)
	return ioutil.WriteFile(f, []byte(s), 0600)
}

func (c *Config) Verify() error {
	if err := validator.Validate(c); err != nil {
		return err
	}

	return nil
}

// DefaultDataDir is the default data directory to use for the databases and other persistence requirements.
func DefaultDataDir() string {
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Application Support", cfgDir)
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", cfgDir)
		} else {
			return filepath.Join(home, nixCfgDir)
		}
	}
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
