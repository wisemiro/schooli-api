package config

import (
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

type AppContext struct {
	Log                     *slog.Logger
	PasswordResetCodeExpiry int
	InProduction            bool
	TokenLifeTime           int
	Wait                    *sync.WaitGroup
	DefaultTimeOut          time.Duration
	Config                  Config
}

func NewAppContext(
	log *slog.Logger,
	passwordResetCodeExpiry int,
	inProduction bool,
	tokenLifeTime int,
	wait *sync.WaitGroup,
	defaultTimeout time.Duration,
	config Config,

) *AppContext {
	return &AppContext{
		Log:                     log,
		PasswordResetCodeExpiry: passwordResetCodeExpiry,
		InProduction:            inProduction,
		TokenLifeTime:           tokenLifeTime,
		Wait:                    wait,
		DefaultTimeOut:          defaultTimeout,
		Config:                  config,
	}
}