//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gophero/goal/aws/s3"
	"github.com/gophero/goal/mailx"
	"github.com/gophero/goal/uuid"
	"goboot/cmd/portal/conf"
	"goboot/cmd/portal/server"
	"goboot/internal/config"
	"goboot/internal/handler"
	"goboot/internal/repository"
	"goboot/internal/repository/repo"
	"goboot/internal/service"

	"github.com/google/wire"
)

// 全局唯一id
var sidSet = wire.NewSet(uuid.NewSid)

var confSet = wire.NewSet(ConfProvider, LoggerProvider, MailConfProvider, wire.Bind(new(config.Conf), new(*conf.PortalConf)))

// http服务
var serverSet = wire.NewSet(server.NewServerHTTP)

// sys相关集合
var sysSrvSet = wire.NewSet(service.NewAppVersionService, service.NewSysSettingService)
var sysRepoSet = wire.NewSet(repository.NewSysSettingRepository, repository.NewSysVersionRepository)

// user相关集合
var userSrvSet = wire.NewSet(service.NewSysDeviceService, service.NewUserService)
var userRepoSet = wire.NewSet(repository.NewSysDeviceRepository, repository.NewUserRepository)

// 集合组
var handlerGroup = wire.NewSet(
	handler.NewHandler, // 顶层handler
)

var srvGroup = wire.NewSet(
	conf.NewS3Conf,
	s3.NewS3,
	mailx.NewSender,
	service.NewService, // 顶层service
	service.NewImageService,
	sysSrvSet,
	userSrvSet,
)

var repoGroup = wire.NewSet(
	repo.NewWDB,
	repo.NewRDB,
	repo.NewRedis,
	repo.NewMongoDB,
	repo.NewRepository,
	sysRepoSet,
	userRepoSet,
)

func newApp() (*server.App, func(), error) {
	panic(wire.Build(
		confSet,
		//sidSet,
		repoGroup,
		//srvGroup,
		//handlerGroup,
		serverSet,
	))
}
