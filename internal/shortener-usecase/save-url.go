package shortener

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strings"

	"github.com/Util787/url-shortener/internal/common"
	"github.com/google/uuid"
)

func (s *shortenerUsecase) SaveURL(longURL string) error {
	ok, _ := validateURL(longURL)
	if !ok {
		return fmt.Errorf("invalid URL: %s", longURL)
	}

	id, shortURL := generateShortURL()

	//validation for duplicates
	exists, err := s.storage.LongURLExists(context.Background(), longURL)
	if err != nil {
		return fmt.Errorf("failed to check if long URL exists: %w", err)
	}
	if exists {
		return fmt.Errorf("long URL already exists: %s", longURL)
	}

	exists, err = s.storage.ShortURLExists(context.Background(), shortURL)
	if err != nil {
		return fmt.Errorf("failed to check if short URL exists: %w", err)
	}
	if exists {
		return fmt.Errorf("short URL already exists: %s", shortURL)
	}

	return s.storage.SaveURL(context.Background(), id, longURL, shortURL)
}

// returns bool and protocol(scheme)
func validateURL(rawURL string) (bool, string) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false, ""
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false, ""
	}

	for _, r := range parsedURL.Scheme {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '+' || r == '-' || r == '.') {
			return false, ""
		}
	}

	return true, strings.ToLower(parsedURL.Scheme)
}

func generateShortURL() (id string, shortURL string) {
	u := uuid.New()

	id = base62Encode(new(big.Int).SetBytes(u[:]))

	shortURL = common.RedirectBaseURL + id
	return id, shortURL
}

func base62Encode(num *big.Int) string {
	const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result []byte

	n := new(big.Int).Set(num)
	base := big.NewInt(62)
	remainder := new(big.Int)

	for n.Sign() > 0 {
		n.DivMod(n, base, remainder)
		result = append([]byte{base62[remainder.Int64()]}, result...)
	}

	return string(result)
}
