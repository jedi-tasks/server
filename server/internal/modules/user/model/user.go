package model

import (
	"github.com/google/uuid"

	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jedi-tasks/go-server-grpc/proto"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`                  // Идентификатор пользователя
	Name         string    `db:"name" json:"name"`              // Имя пользователя
	Email        string    `db:"email" json:"email"`            // Электронная почта
	PasswordHash []byte    `db:"password_hash" json:"-"`        // Хэш пароля
	PasswordSalt []byte    `db:"password_salt" json:"-"`        // Соль пароля
	TimeCreate   time.Time `db:"time_create" json:"timeCreate"` // Дата и время создания аккаунта
}

// ConvertToProto converts User to proto format
func (u User) ConvertToProto() (*proto.User, error) {
	return &proto.User{
		Id:         u.ID[:],
		Name:       u.Name,
		Email:      u.Email,
		TimeCreate: timestamppb.New(u.TimeCreate),
	}, nil
}
