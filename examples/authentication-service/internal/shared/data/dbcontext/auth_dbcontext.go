package dbcontext

import (
	"github.com/phatnt199/go-infra/pkg/infra/postgres/gorm/contracts"
	"github.com/phatnt199/go-infra/pkg/infra/postgres/gorm/dbcontext"
	"gorm.io/gorm"
)

type AuthGormDBContext struct {
	contracts.GormDBContext
}

func NewAuthGormDBContext(db *gorm.DB) AuthGormDBContext {
	c := &AuthGormDBContext{GormDBContext: dbcontext.NewGormDBContext(db)}

	return *c
}
