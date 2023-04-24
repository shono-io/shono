package reaktors

import (
	"fmt"
)

func ErrStorageGetFailed(err error) error {
	return fmt.Errorf("retrieval by id from storage failed: %w", err)
}

func ErrStorageQueryFailed(err error) error {
	return fmt.Errorf("retrieval by query from storage failed: %w", err)
}
