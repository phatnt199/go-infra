package config

import (
	"strings"

	"github.com/iancoleman/strcase"
	appConfig "github.com/phatnt199/go-infra/pkg/application/config"
	"github.com/phatnt199/go-infra/pkg/application/environment"
	"github.com/phatnt199/go-infra/pkg/reflection/typemapper"
)

type MigrationOptions struct {
	Host          string `mapstructure:"host" env:"MIGRATION_HOST" default:"localhost"`
	Port          int    `mapstructure:"port" env:"MIGRATION_PORT" default:"54100"`
	User          string `mapstructure:"user" env:"MIGRATION_USER" default:"admin"`
	Password      string `mapstructure:"password" env:"MIGRATION_PASSWORD" default:"123456"`
	DbName        string `mapstructure:"dbName" env:"MIGRATION_DBNAME" default:"auth_db"`
	SSLMode       bool   `mapstructure:"sslMode" env:"MIGRATION_SSLMODE" default:"false"`
	MigrationsDir string `mapstructure:"migrationsDir" env:"MIGRATION_DIR" default:"db/migrations"`
	SkipMigration bool   `mapstructure:"skipMigration" env:"SKIP_MIGRATION" default:"false"`
}

func NewMigrationOptions(environment environment.Environment) (*MigrationOptions, error) {
	optionName := strcase.ToLowerCamel(typemapper.GetGenericTypeNameByT[MigrationOptions]())

	cfg, err := appConfig.BindConfigKey[*MigrationOptions](optionName, environment)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (o *MigrationOptions) String() string {
	password := o.Password
	if password != "" {
		password = "***"
	}
	return strings.Join([]string{o.User, password, "@", o.Host}, "")
}
