package main

import (
	"context"
	"net"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	ingestrpcdelivry "github.com/osmosis-labs/sqs/ingest/delivery/grpc"
	ingestusecase "github.com/osmosis-labs/sqs/ingest/usecase"

	chaininforepo "github.com/osmosis-labs/sqs/chaininfo/repository"
	chaininfousecase "github.com/osmosis-labs/sqs/chaininfo/usecase"
	poolsHttpDelivery "github.com/osmosis-labs/sqs/pools/delivery/http"
	poolsUseCase "github.com/osmosis-labs/sqs/pools/usecase"
	routerrepo "github.com/osmosis-labs/sqs/router/repository"
	tokenshttpdelivery "github.com/osmosis-labs/sqs/tokens/delivery/http"
	tokensUseCase "github.com/osmosis-labs/sqs/tokens/usecase"

	"github.com/osmosis-labs/sqs/domain"
	"github.com/osmosis-labs/sqs/domain/cache"
	"github.com/osmosis-labs/sqs/domain/mvc"
	"github.com/osmosis-labs/sqs/log"
	"github.com/osmosis-labs/sqs/middleware"

	routerHttpDelivery "github.com/osmosis-labs/sqs/router/delivery/http"
	routerUseCase "github.com/osmosis-labs/sqs/router/usecase"

	systemhttpdelivery "github.com/osmosis-labs/sqs/system/delivery/http"
)

// SideCarQueryServer defines an interface for sidecar query server (SQS).
// It encapsulates all logic for ingesting chain data into the server
// and exposes endpoints for querying formatter and processed data from frontend.
type SideCarQueryServer interface {
	GetRouterRepository() routerrepo.RouterRepository
	GetTokensUseCase() mvc.TokensUsecase
	GetLogger() log.Logger
	Shutdown(context.Context) error
	Start(context.Context) error
}

type sideCarQueryServer struct {
	routerRepository routerrepo.RouterRepository
	tokensUseCase    mvc.TokensUsecase
	e                *echo.Echo
	sqsAddress       string
	logger           log.Logger
}

// GetTokensUseCase implements SideCarQueryServer.
func (sqs *sideCarQueryServer) GetTokensUseCase() mvc.TokensUsecase {
	return sqs.tokensUseCase
}

// GetRouterRepository implements SideCarQueryServer.
func (sqs *sideCarQueryServer) GetRouterRepository() routerrepo.RouterRepository {
	return sqs.routerRepository
}

// GetLogger implements SideCarQueryServer.
func (sqs *sideCarQueryServer) GetLogger() log.Logger {
	return sqs.logger
}

// Shutdown implements SideCarQueryServer.
func (sqs *sideCarQueryServer) Shutdown(ctx context.Context) error {
	return sqs.e.Shutdown(ctx)
}

// Start implements SideCarQueryServer.
func (sqs *sideCarQueryServer) Start(context.Context) error {
	sqs.logger.Info("Starting sidecar query server", zap.String("address", sqs.sqsAddress))
	err := sqs.e.Start(sqs.sqsAddress)
	if err != nil {
		return err
	}

	return nil
}

// NewSideCarQueryServer creates a new sidecar query server (SQS).
func NewSideCarQueryServer(appCodec codec.Codec, config domain.Config, logger log.Logger) (SideCarQueryServer, error) {
	// Setup echo server
	e := echo.New()
	middleware := middleware.InitMiddleware()
	e.Use(middleware.CORS)
	e.Use(middleware.InstrumentMiddleware)
	e.Use(middleware.TraceWithParamsMiddleware("sqs"))

	// Initialize pools repository, usecase and HTTP handler
	poolsUseCase := poolsUseCase.NewPoolsUsecase(config.Pools, config.ChainGRPCGatewayEndpoint)

	// Initialize router repository, usecase
	routerRepository := routerrepo.New()
	routerUsecase := routerUseCase.NewRouterUsecase(routerRepository, poolsUseCase, *config.Router, poolsUseCase.GetCosmWasmPoolConfig(), logger, cache.New(), cache.New())

	// Initialize system handler
	chainInfoRepository := chaininforepo.New()
	chainInfoUseCase := chaininfousecase.NewChainInfoUsecase(chainInfoRepository)

	// Compute token metadata from chain denom.
	tokenMetadataByChainDenom, err := tokensUseCase.GetTokensFromChainRegistry(config.ChainRegistryAssetsFileURL)
	if err != nil {
		return nil, err
	}

	// Initialized tokens usecase
	tokensUseCase := tokensUseCase.NewTokensUsecase(tokenMetadataByChainDenom)

	// HTTP handlers
	poolsHttpDelivery.NewPoolsHandler(e, poolsUseCase)
	systemhttpdelivery.NewSystemHandler(e, config, logger, chainInfoUseCase)
	if err := tokenshttpdelivery.NewTokensHandler(e, *config.Pricing, tokensUseCase, routerUsecase, logger); err != nil {
		return nil, err
	}
	routerHttpDelivery.NewRouterHandler(e, routerUsecase, tokensUseCase, logger)

	// Start grpc ingest server if enabled
	grpcIngesterConfig := config.GRPCIngester
	if grpcIngesterConfig.Enabeld {
		// Initialize ingest handler and usecase
		ingestUseCase, err := ingestusecase.NewIngestUsecase(poolsUseCase, routerUsecase, chainInfoUseCase, tokensUseCase, appCodec, *config.Pricing, logger)
		if err != nil {
			return nil, err
		}

		grpcIngestHandler, err := ingestrpcdelivry.NewIngestGRPCHandler(ingestUseCase, *grpcIngesterConfig)
		if err != nil {
			panic(err)
		}

		go func() {
			logger.Info("Starting grpc ingest server")

			lis, err := net.Listen("tcp", grpcIngesterConfig.ServerAddress)
			if err != nil {
				panic(err)
			}
			if err := grpcIngestHandler.Serve(lis); err != nil {
				panic(err)
			}
		}()
	}

	go func() {
		logger.Info("Starting profiling server")
		err = http.ListenAndServe("localhost:6062", nil)
		if err != nil {
			panic(err)
		}
	}()

	return &sideCarQueryServer{
		routerRepository: routerRepository,
		tokensUseCase:    tokensUseCase,
		logger:           logger,
		e:                e,
		sqsAddress:       config.ServerAddress,
	}, nil
}
