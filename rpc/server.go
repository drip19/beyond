package rpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/drip/beyond/config"
	"github.com/drip/beyond/pkg/log"
	"github.com/drip/beyond/pkg/util"
	"github.com/drip/beyond/rpc/apis"
	pb "github.com/drip/beyond/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Server struct {
	rpc    *grpc.Server
	srv    *http.Server
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func NewServer(cfg *config.Config) *Server {
	gRpcServer := grpc.NewServer(grpc.StreamInterceptor(streamServerInterceptor),
		grpc.UnaryInterceptor(unaryServerInterceptor))

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		cfg:    cfg,
		rpc:    gRpcServer,
		ctx:    ctx,
		cancel: cancel,
		logger: log.NewLogger("rpc"),
	}
}

func (g *Server) Start() error {
	network, address, err := util.Scheme(g.cfg.RPCCfg.GRPCListenAddress)
	if err != nil {
		return err
	}

	lis, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}
	if err := g.registerApi(); err != nil {
		g.logger.Error(err)
		return fmt.Errorf("registerApi: %s", err)
	}
	reflection.Register(g.rpc)
	go func() {
		if err := g.rpc.Serve(lis); err != nil {
			g.logger.Error(err)
		}
	}()
	go func() {
		if err := g.newGateway(address, g.cfg.RPCCfg.ListenAddress); err != nil {
			g.logger.Errorf("gateway listen err: %s", err)
		}
	}()

	g.logger.Info("rpc server started")

	return nil
}

func (g *Server) newGateway(grpcAddress, gwAddress string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	// no need proxy for internal gateway to internal rpc server
	optDial := grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		network := "tcp"
		g.logger.Debugf("WithContextDialer addr %s", addr)
		return (&net.Dialer{}).DialContext(ctx, network, addr)
	})
	opts := []grpc.DialOption{grpc.WithInsecure(), optDial}
	if err := registerGWApi(ctx, gwmux, grpcAddress, opts); err != nil {
		return fmt.Errorf("gateway register: %s", err)
	}
	_, address, err := util.Scheme(gwAddress)
	if err != nil {
		return err
	}
	handler := newCorsHandler(gwmux, g.cfg.RPCCfg.CORSAllowedOrigins)

	g.srv = &http.Server{
		Addr:    address,
		Handler: handler,
	}

	g.srv.RegisterOnShutdown(func() {
		g.logger.Debug("RESEful server shutdown")
	})

	if err := g.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (g *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if g.srv != nil {
		if err := g.srv.Shutdown(ctx); err != nil {
			g.logger.Errorf("RESTful server shutdown failed:%+v", err)
		}
	}
	g.rpc.Stop()
	g.logger.Info("rpc stopped")
}

func (g *Server) registerApi() error {
	pb.RegisterPingAPIServer(g.rpc, apis.NewPingApi(g.cfg))
	return nil
}

func registerGWApi(ctx context.Context, gwmux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := pb.RegisterPingAPIHandlerFromEndpoint(ctx, gwmux, endpoint, opts); err != nil {
		return err
	}
	return nil
}

func newCorsHandler(srv http.Handler, allowedOrigins []string) http.Handler {
	if len(allowedOrigins) == 0 {
		return srv
	}
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{http.MethodPost, http.MethodGet},
		MaxAge:         600,
		AllowedHeaders: []string{"*"},
	})
	return c.Handler(srv)
}

func unaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func streamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	return err
}
