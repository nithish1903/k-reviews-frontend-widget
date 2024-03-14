package elastic

import (
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/matryer/resync"
	"go.uber.org/zap"
)

var onceElastic resync.Once
var elasticClient *elasticsearch.Client

type ElasticCon struct {
	ElasticConn *elasticsearch.Client
}

func ConnectElasticsearch() *ElasticCon {

	onceElastic.Do(func() {
		// Create a new Elasticsearch client
		client, err := elasticsearch.NewClient(elasticsearch.Config{
			CloudID:  os.Getenv("ELASTIC_CLOUD_URL"),
			Username: os.Getenv("ELASTIC_USERNAME"),
			Password: os.Getenv("ELASTIC_PASSWORD"),
		})
		if err != nil {
			zap.L().Fatal("Could not connect to ES", zap.Any("err", err))
		}
		elasticClient = client
		zap.L().Info("Connected to elastic search!")
	})

	return &ElasticCon{ElasticConn: elasticClient}
}

func CreateESConnection() *ElasticCon {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		CloudID:  os.Getenv("ELASTIC_CLOUD_URL"),
		Username: os.Getenv("ELASTIC_USERNAME"),
		Password: os.Getenv("ELASTIC_PASSWORD"),
	})
	if err != nil {
		zap.L().Fatal("Could not connect to ES", zap.Any("err", err))
	}
	elasticClient = client
	zap.L().Info("Connected to elastic search!")

	return &ElasticCon{ElasticConn: elasticClient}
}
