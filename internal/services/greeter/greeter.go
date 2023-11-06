package greeter

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-starterkit/internal/contracts/config"
	fluffycore_contracts_somedisposable "github.com/fluffy-bunny/fluffycore-starterkit/internal/contracts/somedisposable"
	proto_helloworld "github.com/fluffy-bunny/fluffycore-starterkit/proto/helloworld"
	endpoint "github.com/fluffy-bunny/fluffycore/contracts/endpoint"
	grpc_gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	zerolog "github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type (
	service struct {
		proto_helloworld.UnimplementedGreeterServer
		config               *contracts_config.Config
		scopedSomeDisposable fluffycore_contracts_somedisposable.IScopedSomeDisposable
	}
)

func (s *service) SayHello(ctx context.Context, request *proto_helloworld.HelloRequest) (*proto_helloworld.HelloReply, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("SayHello")
	return &proto_helloworld.HelloReply{
		Message: "Hello " + request.Name,
	}, nil
}

type registrationServer struct {
	proto_helloworld.FluffyCoreGreeterServer
}

var _ endpoint.IEndpointRegistration = (*registrationServer)(nil)

func (s *registrationServer) RegisterHandler(gwmux *grpc_gateway_runtime.ServeMux, conn *grpc.ClientConn) {
	proto_helloworld.RegisterGreeterHandler(context.Background(), gwmux, conn)
}
func AddGreeterService(builder di.ContainerBuilder) {
	proto_helloworld.AddGreeterServer[proto_helloworld.IGreeterServer](builder,
		func(config *contracts_config.Config, scopedSomeDisposable fluffycore_contracts_somedisposable.IScopedSomeDisposable) proto_helloworld.IGreeterServer {
			return &service{
				config:               config,
				scopedSomeDisposable: scopedSomeDisposable,
			}
		}, func() endpoint.IEndpointRegistration {
			return &registrationServer{}
		})
}
