package job

import (
	"context"
	"encoding/json"
	"errors"
	"review-job/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
)

//评价数据流处理

// JobWorker 自定义执行job结构体，实现transport.Server接口
type JobWorker struct {
	kafkaReader *kafka.Reader
	esClient    *ESClient
	log         *log.Helper
}

type ESClient struct {
	*elasticsearch.TypedClient
	index string
}

func NewJobWorker(kafkaReader *kafka.Reader, esClient *ESClient, logger log.Logger) *JobWorker {
	return &JobWorker{
		kafkaReader: kafkaReader,
		esClient:    esClient,
		log:         log.NewHelper(logger),
	}
}

func NewKafkaReader(cfg *conf.Kafka) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GropuId,
		Topic:   cfg.Topic,
	})
}

func NewESClient(cfg *conf.Elasticsearch) (*ESClient, error) {
	//ES 配置
	c := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}
	//创建ES客户端
	client, err := elasticsearch.NewTypedClient(c)
	if err != nil {
		return nil, err
	}
	return &ESClient{
		TypedClient: client,
		index:       cfg.Index,
	}, nil
}

// Msg定义kafka中接收到的数据
type Msg struct {
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	IsDdl    bool                     `json:"isDdl"`
	Data     []map[string]interface{} `json:"data"`
}

// Start程序启动之后干活的
// ctx 是kratos 启动时传入的ctx,是带有取消的
func (jw JobWorker) Start(ctx context.Context) error {
	jw.log.Debug("JobWorker Start.......")
	// 1.从kafka中获取MYSQL中的数据变更消息
	// 接收消息
	for {
		m, err := jw.kafkaReader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) { // 如果是ctx取消的错误
			return nil
		}
		if err != nil {
			jw.log.Error("kafka.Reader.ReadMessage failed, err:%v", err)
			break
		}
		jw.log.Debugf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 2.将完整评价数据写入ES中
		msg := new(Msg)
		err = json.Unmarshal(m.Value, msg) // 解析消息到msg中
		if err != nil {
			jw.log.Error("json.Unmarshal failed, err:%v", err)
			continue
		}

		//对数据进行业务处理

		if msg.Type == "INSERT" {
			//往ES中写入文档
			for idx := range msg.Data {
				jw.indexDocument(msg.Data[idx])
			}
		} else {
			//往ES中更新文档
			for idx := range msg.Data {
				jw.updateDocument(msg.Data[idx])
			}
		}

	}

	return nil
}

// Stop kratos 程序启动之后会调用的方法
func (jw JobWorker) Stop(ctx context.Context) error {
	jw.log.Debug("JobWorker Stop.......")
	// 程序退出前关闭Reader
	return jw.kafkaReader.Close()
}
func (jw JobWorker) connkafka() {
	jw.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		GroupID:  "consumer-group-id", // 指定消费者组id
		Topic:    "topic-A",
		MaxBytes: 10e6, // 10MB
	})
}

// indexDocument 索引文档
func (jw JobWorker) indexDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)
	// 添加文档
	resp, err := jw.esClient.Index(jw.esClient.index).
		Id(reviewID).
		Document(d).
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("indexing document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%#v\n", resp.Result)
}

// updateDocument 更新文档
func (jw JobWorker) updateDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)
	resp, err := jw.esClient.Update(jw.esClient.index, reviewID).
		Doc(d). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("update document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%v\n", resp.Result)
}
