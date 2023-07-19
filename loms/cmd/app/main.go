package main

import (
	"context"
	"net"
	"net/http"
	api "route256/loms/internal/api/loms"
	"route256/loms/internal/config"
	"route256/loms/internal/kafka"
	"route256/loms/internal/model"
	"route256/loms/internal/pkg/logger"
	"route256/loms/internal/pkg/tracer"
	"route256/loms/internal/repository/postgres"
	"route256/loms/internal/repository/postgres/tx"
	"route256/loms/internal/service/loms"
	"route256/loms/pkg/loms_v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":50052"

func main() {
	if err := tracer.InitGlobal(model.ServiceName); err != nil {
		logger.Fatal("ERR: ", err)
	}

	err := config.Init()
	if err != nil {
		logger.Fatal("ERR: ", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, config.AppConfig.PostgresDB.URL)
	if err != nil {
		logger.Fatal("connect to db: %w", err)
	}
	defer pool.Close()

	provider := tx.New(pool)
	repo := postgres.NewRepo(provider)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		mux := runtime.NewServeMux(
			runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
				switch key {
				case "x-trace-id":
					return key, true
				}
				return runtime.DefaultHeaderMatcher(key)
			}),
		)

		if err := mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			promhttp.Handler().ServeHTTP(w, r)
		}); err != nil {
			logger.Fatal("something wrong with metrics handler", err)
		}

		if err := loms_v1.RegisterLomsHandlerFromEndpoint(gCtx, mux, ":50052", []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}); err != nil {
			return errors.Wrap(err, "cannot register http server")
		}

		httpHost := ":8081"
		logger.Info("HTTP server started on: ", httpHost)

		if err := http.ListenAndServe(httpHost, mux); err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			return err
			// logger.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				logger.MiddlewareGRPC,
				tracer.MiddlewareGRPC,
			),
		)

		reflection.Register(s)

		kafkaProducer, err := kafka.NewProducer(config.AppConfig.Kafka.Brokers, config.AppConfig.Kafka.Topic)
		if err != nil {
			logger.Fatalf("failed to kafka produce: %v", err)
		}

		loms_v1.RegisterLomsServer(s, api.NewLomsServer(loms.NewService(provider, repo, kafkaProducer)))

		logger.Info("server listening at %v", lis.Addr())

		if err = s.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
		return nil
	})

	err = g.Wait()
	if err != nil {
		logger.Fatal("g wait")
	}
}
