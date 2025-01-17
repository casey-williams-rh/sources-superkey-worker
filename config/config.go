package config

import (
	"fmt"
	"os"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
)

// SuperKeyWorkerConfig is the struct for storing runtime configuration
type SuperKeyWorkerConfig struct {
	Hostname           string
	KafkaBrokers       []string
	KafkaTopics        map[string]string
	KafkaGroupID       string
	MetricsPort        int
	LogLevel           string
	LogGroup           string
	LogHandler         string
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
	SourcesHost        string
	SourcesScheme      string
	SourcesPort        int
	SourcesPSK         string
}

// Get - returns the config parsed from runtime vars
func Get() *SuperKeyWorkerConfig {
	options := viper.New()
	kafkaTopics := make(map[string]string)

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		for requestedName, topicConfig := range clowder.KafkaTopics {
			kafkaTopics[requestedName] = topicConfig.Name
		}
		options.SetDefault("AwsRegion", cfg.Logging.Cloudwatch.Region)
		options.SetDefault("AwsAccessKeyId", cfg.Logging.Cloudwatch.AccessKeyId)
		options.SetDefault("AwsSecretAccessKey", cfg.Logging.Cloudwatch.SecretAccessKey)
		options.SetDefault("KafkaBrokers", []string{fmt.Sprintf("%s:%v", cfg.Kafka.Brokers[0].Hostname, *cfg.Kafka.Brokers[0].Port)})
		options.SetDefault("LogGroup", cfg.Logging.Cloudwatch.LogGroup)
		options.SetDefault("MetricsPort", cfg.MetricsPort)

	} else {
		options.SetDefault("AwsRegion", "us-east-1")
		options.SetDefault("AwsAccessKeyId", os.Getenv("CW_AWS_ACCESS_KEY_ID"))
		options.SetDefault("AwsSecretAccessKey", os.Getenv("CW_AWS_SECRET_ACCESS_KEY"))
		options.SetDefault("KafkaBrokers", []string{fmt.Sprintf("%v:%v", os.Getenv("QUEUE_HOST"), os.Getenv("QUEUE_PORT"))})
		options.SetDefault("LogGroup", os.Getenv("CLOUD_WATCH_LOG_GROUP"))
		options.SetDefault("MetricsPort", 9394)

	}

	options.SetDefault("KafkaGroupID", "sources-superkey-worker")
	options.SetDefault("KafkaTopics", kafkaTopics)
	options.SetDefault("LogLevel", os.Getenv("LOG_LEVEL"))
	options.SetDefault("LogHandler", os.Getenv("LOG_HANDLER"))

	options.SetDefault("SourcesHost", os.Getenv("SOURCES_HOST"))
	options.SetDefault("SourcesScheme", os.Getenv("SOURCES_SCHEME"))
	options.SetDefault("SourcesPort", os.Getenv("SOURCES_PORT"))
	options.SetDefault("SourcesPSK", os.Getenv("SOURCES_PSK"))

	hostname, _ := os.Hostname()
	options.SetDefault("Hostname", hostname)

	options.AutomaticEnv()

	return &SuperKeyWorkerConfig{
		Hostname:           options.GetString("Hostname"),
		KafkaBrokers:       options.GetStringSlice("KafkaBrokers"),
		KafkaTopics:        options.GetStringMapString("KafkaTopics"),
		KafkaGroupID:       options.GetString("KafkaGroupID"),
		MetricsPort:        options.GetInt("MetricsPort"),
		LogLevel:           options.GetString("LogLevel"),
		LogHandler:         options.GetString("LogHandler"),
		LogGroup:           options.GetString("LogGroup"),
		AwsRegion:          options.GetString("AwsRegion"),
		AwsAccessKeyID:     options.GetString("AwsAccessKeyID"),
		AwsSecretAccessKey: options.GetString("AwsSecretAccessKey"),
		SourcesHost:        options.GetString("SourcesHost"),
		SourcesScheme:      options.GetString("SourcesScheme"),
		SourcesPort:        options.GetInt("SourcesPort"),
		SourcesPSK:         options.GetString("SourcesPSK"),
	}
}
