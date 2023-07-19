package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"schooli-api/cmd/config"
	"schooli-api/internal/auth/token"
	"schooli-api/internal/handlers"
	"schooli-api/internal/services"
	"schooli-api/pkg/cache"
	"schooli-api/pkg/filestore"
	"schooli-api/pkg/logger"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/hibiken/asynq"
	"golang.org/x/exp/slog"
)

type Application struct {
	logger      *slog.Logger
	mux         *chi.Mux
	wg          *sync.WaitGroup
	fileStorage filestore.FileStorage
	cache       cache.Store
	appCtx      *config.AppContext
	redisOpt    asynq.RedisClientOpt
	conf        config.Config
	store       services.Store
}

func StartApplication(loc string) (*Application, error) {
	l := logger.GetLogger()
	conf, err := config.Load(loc)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}

	// Get file storage connection
	redisCache, err := cache.NewCache(conf.Redis.Address, conf.Redis.Port, l)
	if err != nil {
		return nil, err
	}
	appContext := config.NewAppContext(
		l,
		20,
		conf.Server.Prod,
		20,
		wg,
		time.Duration(conf.Server.Timeout)*time.Second,
		conf,
	)
	// initiate services
	tokenMaker, err := token.NewTokenMaker(conf.Secrets.Jwt, appContext)
	if err != nil {
		return nil, err

	}
	// Get file storage connection
	s3Storage, err := filestore.NewS3Storage(filestore.FileBackendSettings{
		AmazonS3AccessKeyId:                conf.Minio.User,
		AmazonS3SecretAccessKey:            conf.Minio.Password,
		AmazonS3Bucket:                     conf.Minio.BucketName,
		AmazonS3PathPrefix:                 "",
		AmazonS3Region:                     conf.Minio.Region,
		AmazonS3Endpoint:                   conf.Minio.Address,
		AmazonS3SSL:                        false,
		AmazonS3SignV2:                     false,
		AmazonS3SSE:                        false,
		AmazonS3Trace:                      false,
		SkipVerify:                         true,
		AmazonS3RequestTimeoutMilliseconds: 10,
	})
	if err != nil {
		return nil, err
	}

	store, err := services.NewSQLStore(conf, l, s3Storage)
	if err != nil {
		return nil, err
	}
	allServices := handlers.NewRepo(
		s3Storage,
		store,
		tokenMaker,
	)

	mux := allServices.SetupRouter()
	return &Application{
		logger:      l,
		mux:         mux,
		wg:          wg,
		fileStorage: s3Storage,
		cache:       redisCache,
		appCtx:      appContext,
		redisOpt: asynq.RedisClientOpt{
			Addr: fmt.Sprintf("%s:%s", conf.Redis.Address, conf.Redis.Port),
		},
		conf:  conf,
		store: store,
	}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:              a.conf.Server.Address,
		Handler:           a.mux,
		ReadTimeout:       time.Duration(a.conf.Server.Timeout) * time.Second,
		ReadHeaderTimeout: time.Duration(a.conf.Server.Timeout) * time.Second,
		WriteTimeout:      time.Duration(a.conf.Server.Timeout) * time.Second,
		IdleTimeout:       time.Duration(a.conf.Server.Timeout) * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")
	return nil
}
