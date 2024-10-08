// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/gophero/goal/aws/s3"
	"github.com/gophero/goal/mailx"
	"github.com/gophero/goal/uuid"
	"goboot/cmd/job/conf"
	"goboot/cmd/job/server"
	"goboot/cmd/job/task"
	"goboot/internal/config"
	"goboot/internal/repository"
	"goboot/internal/repository/repo"
	"goboot/internal/service"
)

// Injectors from wire.go:

func newApp() (*server.App, func(), error) {
	jobConf := ConfProvider()
	conf := MailConfProvider(jobConf)
	logger := LoggerProvider(jobConf)
	client := repo.NewRedis(jobConf)
	wdb := repo.NewWDB(jobConf)
	rdb := repo.NewRDB(jobConf)
	repository := repo.NewRepository(wdb, rdb, logger)
	sid := uuid.NewSid()
	serviceService := service.NewService(repository, logger, sid, jobConf)
	demoTask := task.NewDemoTask(serviceService, client)
	app := server.NewServerHTTP(conf, logger, jobConf, client, demoTask)
	return app, func() {
	}, nil
}

// wire.go:

var confSet = wire.NewSet(ConfProvider, LoggerProvider, MailConfProvider, wire.Bind(new(config.Conf), new(*conf.JobConf)))

// http服务
var serverSet = wire.NewSet(server.NewServerHTTP)

// Job相关集合
var jobSrvSet = wire.NewSet(task.NewDemoTask)

// 全局唯一id
var sidSet = wire.NewSet(uuid.NewSid)

// sys相关集合
var sysSrvSet = wire.NewSet(service.NewSysSettingService, service.NewSysDeviceService)

var sysRepoSet = wire.NewSet(repository.NewSysSettingRepository, repository.NewSysDeviceRepository)

// user相关集合
var usrSrvSet = wire.NewSet(service.NewUserService)

var usrRepoSet = wire.NewSet(repository.NewUserRepository)

// 集合组
var srvGroup = wire.NewSet(conf.NewS3Conf, s3.NewS3, mailx.NewSender, service.NewService, service.NewImageService, usrSrvSet,
	jobSrvSet,
	sysSrvSet,
)

var repoGroup = wire.NewSet(repo.NewWDB, repo.NewRDB, repo.NewRedis, repo.NewMongoDB, repo.NewRepository, usrRepoSet,
	sysRepoSet,
)
