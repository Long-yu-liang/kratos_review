package data

import (
	"errors"
	"fmt"
	"review-service/internal/conf"
	"review-service/internal/data/query"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewReviewRepo, NewDB, NewESClient, NewRedis)

// Data .
type Data struct {
	// TODO wrapped database client
	query *query.Query
	log   *log.Helper
	es    *elasticsearch.TypedClient
	rdb   *redis.Client
}

// NewData .
func NewData(db *gorm.DB, esClient *elasticsearch.TypedClient, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)

	return &Data{query: query.Q, es: esClient, rdb: rdb, log: log.NewHelper(logger)}, cleanup, nil
}

func NewESClient(cfg *conf.Elasticsearch) (*elasticsearch.TypedClient, error) {
	//ES 配置
	c := elasticsearch.Config{
		Addresses: cfg.GetAddresses(),
	}

	//创建客户端连接
	return elasticsearch.NewTypedClient(c)
}

func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.Database.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db, nil
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(cfg.Database.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db, nil
	}
	return nil, errors.New("connectDB failed unsupport driver")
}

func NewRedis(cfg *conf.Data) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		ReadTimeout:  cfg.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: cfg.Redis.WriteTimeout.AsDuration(),
	})
}
