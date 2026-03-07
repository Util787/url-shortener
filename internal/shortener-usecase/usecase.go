package shortener

import "context"

type shortenerUsecase struct {
	storage         URLMappingStorage
	redirectBaseURL string
}

type URLMappingStorage interface {
	SaveURL(ctx context.Context, id string, longURL string, shortURL string) error
	GetLongURL(ctx context.Context, shortURL string) (string, error)
	GetRandomURL(ctx context.Context) (string, error)
	LongURLExists(ctx context.Context, longURL string) (bool, error)
	ShortURLExists(ctx context.Context, shortURL string) (bool, error)
	DeleteURL(ctx context.Context, id *string, longURL *string, shortURL *string) error
}

func NewShortenerUsecase(storage URLMappingStorage, redirectBaseURL string) *shortenerUsecase {
	return &shortenerUsecase{
		storage:         storage,
		redirectBaseURL: redirectBaseURL,
	}
}
