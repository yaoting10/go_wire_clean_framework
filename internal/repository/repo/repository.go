package repo

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/redisx"
	"goboot/internal/config"
	"io"
	slog "log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// wire 不能使用相同类型的参数，所以需要新建类型

type WDB *gorm.DB
type RDB *gorm.DB

type Repository struct {
	W *gorm.DB
	R *gorm.DB
	L *logx.Logger

	Manager
}

func NewRepository(wdb WDB, rdb RDB, logger *logx.Logger) *Repository {
	return &Repository{
		W:       wdb,
		R:       rdb,
		L:       logger,
		Manager: &GormManager{rdb, wdb},
	}
}

func NewWDB(conf config.Conf) WDB {
	logLevel := logger.Warn
	if conf.Storage().DBsConf.Wdb().ShowSql {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		slog.New(io.MultiWriter(os.Stdout), "\r\n", slog.LstdFlags|slog.Lshortfile), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	db, err := gorm.Open(mysql.Open(conf.Storage().DBsConf.Rdb().Dsn()), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}
	return db
}

func NewRDB(conf config.Conf) RDB {
	logLevel := logger.Warn
	if conf.Storage().DBsConf.Rdb().ShowSql {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		slog.New(io.MultiWriter(os.Stdout), "\r\n", slog.LstdFlags|slog.Lshortfile), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	db, err := gorm.Open(mysql.Open(conf.Storage().DBsConf.Rdb().Dsn()), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}
	return db
}

// NewRedis 文档见：https://redis.uptrace.dev/zh/guide/
func NewRedis(conf config.Conf) redisx.Client {
	var client redis.UniversalClient
	var redisConf = conf.Storage().RedisConf
	var addrs = strings.Split(redisConf.Addr, ",")
	if redisConf.Cluster {
		if redisConf.Tls {
			client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    addrs,
				Password: redisConf.Password,
				//DB:       conf.Redis.DB,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			})
		} else {
			client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    addrs,
				Password: redisConf.Password,
				//DB:       conf.Redis.DB,
			})
		}
		err := client.(*redis.ClusterClient).ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
			r := shard.Ping(ctx)
			s, err := r.Result()
			fmt.Printf("Pinging redis node: %s, result: %s\n", shard.ClientInfo(ctx).Val().LAddr, s)
			return err
		})
		if err != nil {
			panic(err)
		}
	} else {
		if redisConf.Tls {
			client = redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs:    addrs,
				Password: redisConf.Password,
				DB:       redisConf.DB,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			})
		} else {
			client = redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs:    addrs,
				Password: redisConf.Password,
				DB:       redisConf.DB,
			})
		}
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			panic(fmt.Sprintf("Redis error: %s", err.Error()))
		}
	}
	return client
}

func NewMongoDB(conf config.Conf) *mongo.Database {
	var mongodbConf = conf.Storage().MongoConf
	credential := options.Credential{
		AuthSource: mongodbConf.AuthSource,
		Username:   mongodbConf.Username,
		Password:   mongodbConf.Password,
	}
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://" + mongodbConf.Host + ":" + mongodbConf.Port).SetAuth(credential)
	clientOptions.SetMaxPoolSize(10)
	Client, err := mongo.Connect(ctx, clientOptions)
	errorx.Throw(err)
	mdb := Client.Database(mongodbConf.Database)
	return mdb
}
