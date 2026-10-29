package model

import (
	"github.com/google/uuid"

	"server/internal/utils/errors"
)

// parseUUIDs преобразует идентификаторы из байтов в UUID
func parseUUIDs(ids [][]byte) ([]uuid.UUID, error) {
	res := make([]uuid.UUID, 0, len(ids))
	for _, idBytes := range ids {
		id, err := uuid.FromBytes(idBytes)
		if err != nil {
			return nil, errors.BadRequest.Wrap(err)
		}
		res = append(res, id)
	}
	return res, nil
}
