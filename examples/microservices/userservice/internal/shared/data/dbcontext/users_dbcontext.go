package dbcontext

import (
	"github.com/phatnt199/go-infra/pkg/infra/postgres/gorm/contracts"
	"github.com/phatnt199/go-infra/pkg/infra/postgres/gorm/dbcontext"
	"gorm.io/gorm"
)

type UsersGormDBContext struct {
	contracts.GormDBContext
}

func NewUsersGormDBContext(db *gorm.DB) UsersGormDBContext {
	c := &UsersGormDBContext{GormDBContext: dbcontext.NewGormDBContext(db)}

	return *c
}
