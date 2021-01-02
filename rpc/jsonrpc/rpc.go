package jsonrpc

import (
	"errors"
	"github.com/drip/beyond/config"
	"github.com/drip/beyond/pkg/log"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	jsonrpc2 "github.com/drip/beyond/pkg/jsonrpc2"
	"go.uber.org/zap"
)

type RPC struct {
	rpcAPIs          []jsonrpc2.API
	inProcessHandler *jsonrpc2.Server

	ipcListener net.Listener
	ipcHandler  *jsonrpc2.Server

	httpWhitelist []string
	httpListener  net.Listener
	httpHandler   *jsonrpc2.Server

	wsListener net.Listener
	wsHandler  *jsonrpc2.Server

	config *config.Config

	lock   sync.RWMutex
	logger *zap.SugaredLogger
}

func NewRPC(cfg *config.Config) (*RPC, error) {
	r := RPC{
		config: cfg,
		logger: log.NewLogger("grpc"),
	}
	return &r, nil
}

// startIPC initializes and starts the IPC RPC endpoint.
func (r *RPC) startIPC(apis []jsonrpc2.API) error {
	if r.config.RPCCfg.IPCEndpoint == "" {
		return nil // IPC disabled.
	}
	listener, handler, err := jsonrpc2.StartIPCEndpoint(r.config.RPCCfg.IPCEndpoint, apis)
	if err != nil {
		return err
	}
	r.ipcListener = listener
	r.ipcHandler = handler
	r.logger.Info("IPC endpoint opened, ", "url:", r.config.RPCCfg.IPCEndpoint)
	return nil
}

// stopIPC terminates the IPC RPC endpoint.
func (r *RPC) stopIPC() {
	if r.ipcListener != nil {
		r.ipcListener.Close()
		r.ipcListener = nil

		r.logger.Debug("IPC endpoint closed, ", "endpoint:", r.config.RPCCfg.IPCEndpoint)
	}
	if r.ipcHandler != nil {
		r.ipcHandler.Stop()
		r.ipcHandler = nil
	}
}

// startHTTP initializes and starts the HTTP RPC endpoint.
func (r *RPC) startHTTP(endpoint string, apis []jsonrpc2.API, modules []string, cors []string, vhosts []string, timeouts jsonrpc2.HTTPTimeouts, exposeAll bool) error {
	// Short circuit if the HTTP endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := jsonrpc2.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts, timeouts)
	if err != nil {
		return err
	}
	r.logger.Info("HTTP endpoint opened,", " url:", listener.Addr(), ", cors:", strings.Join(cors, ","), ", vhosts:", strings.Join(vhosts, ","))
	// All listeners booted successfully
	//r.httpEndpoint = endpoint
	r.httpListener = listener
	r.httpHandler = handler

	return nil
}

// stopHTTP terminates the HTTP RPC endpoint.
func (r *RPC) stopHTTP() {
	if r.httpListener != nil {
		r.httpListener.Close()
		r.httpListener = nil

		r.logger.Debug("HTTP endpoint closed, ", "endpoint:", r.config.RPCCfg.HTTPEndpoint)
	}
	if r.httpHandler != nil {
		r.httpHandler.Stop()
		r.httpHandler = nil
	}
}

// startWS initializes and starts the websocket RPC endpoint.
func (r *RPC) startWS(endpoint string, apis []jsonrpc2.API, modules []string, wsOrigins []string, exposeAll bool) error {
	// Short circuit if the WS endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := jsonrpc2.StartWSEndpoint(endpoint, apis, modules, wsOrigins, exposeAll)
	if err != nil {
		return err
	}
	r.logger.Info("WebSocket endpoint opened, ", "url:", listener.Addr())
	// All listeners booted successfully
	//r.wsEndpoint = endpoint
	r.wsListener = listener
	r.wsHandler = handler

	return nil
}

// stopWS terminates the websocket RPC endpoint.
func (r *RPC) stopWS() {
	if r.wsListener != nil {
		r.wsListener.Close()
		r.wsListener = nil
		r.logger.Debug("WebSocket endpoint closed, ", "endpoint:", r.config.RPCCfg.WSEndpoint)
	}
	if r.wsHandler != nil {
		r.wsHandler.Stop()
		r.wsHandler = nil
	}
}

func (r *RPC) Attach() (*jsonrpc2.Client, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	//if r.p2pServer == nil {
	//	return nil, ErrNodeStopped
	//}

	if r.inProcessHandler == nil {
		return nil, errors.New("server not started")
	}
	return jsonrpc2.DialInProc(r.inProcessHandler), nil
}

// startInProc initializes an in-process RPC endpoint.
func (r *RPC) startInProcess(apis []jsonrpc2.API) error {
	// Register all the APIs exposed by the services
	handler := jsonrpc2.NewServer()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			r.logger.Info(err)
			return err
		}
		//r.logger.Debug("InProc registered ", "service ", api.Service, " namespace ", api.Namespace)
	}
	r.inProcessHandler = handler
	return nil
}

// stopInProc terminates the in-process RPC endpoint.
func (r *RPC) stopInProcess() {
	if r.inProcessHandler != nil {
		r.inProcessHandler.Stop()
		r.inProcessHandler = nil
	}
}

func (r *RPC) StopRPC() {
	r.stopInProcess()
	if r.config.RPCCfg.Enable && r.config.RPCCfg.IPCEnabled {
		r.stopIPC()
	}
	if r.config.RPCCfg.Enable && r.config.RPCCfg.HTTPEnabled {
		r.stopHTTP()
	}
	if r.config.RPCCfg.Enable && r.config.RPCCfg.WSEnabled {
		r.stopWS()
	}

}

func (r *RPC) StartRPC() error {

	// Init grpc log
	//rpcapi.Init(node.config.DataDir, node.config.LogLevel, node.config.TestTokenHexPrivKey, node.config.TestTokenTti)

	// Start the various API endpoints, terminating all in case of errors
	//if err := r.startInProcess(r.GetInProcessApis()); err != nil {
	//	return err
	//}

	//Start grpc
	if r.config.RPCCfg.Enable && r.config.RPCCfg.IPCEnabled {
		api := r.GetIpcApis()
		if err := r.startIPC(api); err != nil {
			r.stopInProcess()
			return err
		}
	}

	if r.config.RPCCfg.Enable && r.config.RPCCfg.HTTPEnabled {
		apis := r.GetHttpApis()
		timeout := jsonrpc2.HTTPTimeouts{
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 200 * time.Second,
			IdleTimeout:  200 * time.Second,
		}
		if err := r.startHTTP(r.config.RPCCfg.HTTPEndpoint, apis, nil, r.config.RPCCfg.HTTPCors, r.config.RPCCfg.HttpVirtualHosts, timeout, false); err != nil {
			r.logger.Info(err)
			r.stopInProcess()
			r.stopIPC()
			return err
		}
	}

	if r.config.RPCCfg.Enable && r.config.RPCCfg.WSEnabled {
		apis := r.GetWSApis()
		if err := r.startWS(r.config.RPCCfg.WSEndpoint, apis, nil, []string{}, false); err != nil {
			r.logger.Info(err)
			//r.stopInProcess()
			r.stopIPC()
			r.stopHTTP()
			return err
		}
	}
	return nil
}

func scheme(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}
	return u.Scheme, u.Host, nil
}
