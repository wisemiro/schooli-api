package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"schooli-api/assets"
	"schooli-api/cmd/config"
	"schooli-api/internal/repository/postgresql/db"
	"schooli-api/pkg/filestore"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slog"
)

type Store interface {
	UserService
	ProductsService
	OrderService
	RatingsService
}

type SQLStore struct {
	pool      *pgxpool.Pool
	store     db.Store
	fileStore filestore.FileStorage
	log       *slog.Logger
	cfg       config.Config
}

func NewSQLStore(cfg config.Config, log *slog.Logger, storage filestore.FileStorage) (*SQLStore, error) {
	pgpool, err := New(log, cfg)
	if err != nil {
		return nil, err
	}
	store := &SQLStore{
		store:     db.NewStore(pgpool),
		fileStore: storage,
		log:       log,
		cfg:       cfg,
		pool:      pgpool,
	}
	return store, nil
}

// New instantiates the Postgres database using configuration defined in environment variables.
func New(l *slog.Logger, cfg config.Config) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Database.User, cfg.Database.Password),
		Host:   fmt.Sprintf("%s:%v", cfg.Database.Server, cfg.Database.Port),
		Path:   cfg.Database.Database,
	}
	if !cfg.Database.Secure {
		dsn.RawQuery = "sslmode=disable"

	}
	poolConfig, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, err
	}
	poolConfig.HealthCheckPeriod = 10 * time.Second
	if cfg.Database.Log {
		poolConfig.ConnConfig.Tracer = &pgxTracer{
			log: l,
		}
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return conn.Ping(ctx)
	}
	pool, errx := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if errx != nil {
		return nil, errx
	}
	l.Info("Database connected")
	if cfg.Database.Automigrate {
		return pool, validateSchema(pool.Config().ConnString())
	}
	return pool, nil
}

type pgxTracer struct {
	log *slog.Logger
}

// TraceQueryStart will for now implement only logging with the default logger as
// slog  in future will add otel tracing.
func (t *pgxTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// t.log.InfoCtx(ctx, "", slog.Any("query", data.SQL), slog.Any("args", data.Args))
	return ctx
}

func (t *pgxTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	// t.logTracer.TraceQueryEnd(ctx, conn, data)
	// t.finishSpan(ctx, data.Err)
}

func validateSchema(dsn string) error {
	iofsDriver, errx := iofs.New(assets.EmbeddedFiles, "migrations")
	if errx != nil {
		return errx
	}

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error("unable to close db", err)
		}
	}(database)
	targetInstance, err := postgres.WithInstance(database, new(postgres.Config))
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithInstance("iofs", iofsDriver, "postgres", targetInstance)
	if err != nil {
		return err
	}
	err = migrator.Up()
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		break
	case err != nil:
		return err
	}
	slog.Info("Migrations Done")
	return nil
}

func (s *SQLStore) Close() {
	slog.Info("Closing pgx pool", nil)
	s.pool.Close()
}
