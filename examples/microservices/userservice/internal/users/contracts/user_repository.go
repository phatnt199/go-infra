package contracts

import (
	"context"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.User], error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
