package handlers

import (
	"schooli-api/internal/auth/token"
	"schooli-api/internal/services"
	"schooli-api/pkg/filestore"
)

type Repository struct {
	store          services.Store
	storageService filestore.FileStorage
	tokenMaker     token.Maker
}

// NewRepo creates a new repository.
func NewRepo(
	storageService filestore.FileStorage,
	store services.Store,
	tokenMaker token.Maker,

) *Repository {
	return &Repository{
		storageService: storageService,
		store:          store,
		tokenMaker:     tokenMaker,
	}
}
