package model

import "github.com/google/uuid"

type UpdateUserReq struct {
	ID           uuid.UUID
	Name         *string
	Email        *string
	PasswordHash *[]byte
	PasswordSalt *[]byte
}
