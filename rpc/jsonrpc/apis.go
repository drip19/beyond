package jsonrpc

import (
	"github.com/drip/beyond/rpc/jsonrpc/api"
	jsonrpc2 "github.com/drip/beyond/pkg/jsonrpc2"
)

func (r *RPC) getApi(apiModule string) jsonrpc2.API {
	switch apiModule {
	case "ping":
		return jsonrpc2.API{
			Namespace: "ping",
			Version:   "1.0",
			Service:   api.NewPingApi(r.config),
			Public:    true,
		}
	default:
		return jsonrpc2.API{}
	}
}

func (r *RPC) GetApis(apiModule ...string) []jsonrpc2.API {
	var apis []jsonrpc2.API
	for _, m := range apiModule {
		apis = append(apis, r.getApi(m))
	}
	return apis
}

//In-proc apis
func (r *RPC) GetInProcessApis() []jsonrpc2.API {
	return r.GetPublicApis()
}

//Ipc apis
func (r *RPC) GetIpcApis() []jsonrpc2.API {
	return r.GetPublicApis()
}

//Http apis
func (r *RPC) GetHttpApis() []jsonrpc2.API {
	return r.GetPublicApis()
}

//WS apis
func (r *RPC) GetWSApis() []jsonrpc2.API {
	return r.GetPublicApis()
}

func (r *RPC) GetPublicApis() []jsonrpc2.API {
	apiModules := []string{"ping"}
	return r.GetApis(apiModules...)
}
