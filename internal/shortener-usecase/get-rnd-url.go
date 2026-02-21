package shortener

import "context"

func (s *shortenerUsecase) GetRandomURL() (string, error) {
	return s.storage.GetRandomURL(context.Background())
}
