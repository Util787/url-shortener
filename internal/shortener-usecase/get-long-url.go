package shortener

import "context"

func (s *shortenerUsecase) GetLongURL(shortURL string) (string, error) {
	longURL, err := s.storage.GetLongURL(context.Background(), shortURL)
	if err != nil {
		return "", err
	}

	return longURL, nil
}
