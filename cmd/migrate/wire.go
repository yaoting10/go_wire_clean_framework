//go:build wireinject
// +build wireinject

package main

import (
	"goboot/cmd/migrate/conf"
	"goboot/internal/config"
	"goboot/internal/migration"
	"goboot/internal/repository/repo"

	"github.com/google/wire"
)

var confSet = wire.NewSet(ConfProvider, LoggerProvider, wire.Bind(new(config.Conf), new(*conf.MigrateConf)))

var repoSet = wire.NewSet(
	repo.NewWDB,
	repo.NewRepository,
)
var migrateSet = wire.NewSet(migration.NewMigrate)

func newApp() (*migration.Migrate, func(), error) {
	panic(wire.Build(
		confSet,
		repoSet,
		migrateSet,
	))
}
