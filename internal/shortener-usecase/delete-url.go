package shortener

import (
	"context"
	"fmt"

	"github.com/Util787/url-shortener/internal/common"
)

func (s *shortenerUsecase) DeleteURL(id *string, longURL *string, shortURL *string) error {
	op := common.GetOperationName()

	if id == nil && longURL == nil && shortURL == nil {
		return fmt.Errorf("%s: no criteria for delete provided", op) // No criteria provided, nothing to delete
	}
	return s.storage.DeleteURL(context.Background(), id, longURL, shortURL)
}
