package queues

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	AwsKey []byte
	AwsSecret []byte
	AwsRegion []byte
	GdprQueueUrl []byte
	GdprCallbackQueueUrl []byte
	GdprBucket []byte
	GdprBasePath []byte
}

func InitConfig() (*Config, error) {
	config := &Config{
		AwsKey: []byte(viper.GetString("aws_key")),
		AwsSecret: []byte(viper.GetString("aws_secret")),
		AwsRegion: []byte(viper.GetString("aws_region")),
		GdprQueueUrl: []byte(viper.GetString("gdpr_queue_url")),
		GdprCallbackQueueUrl: []byte(viper.GetString("gdpr_callback_queue_url")),
		GdprBucket: []byte(viper.GetString("gdpr_bucket")),
		GdprBasePath: []byte(viper.GetString("gdpr_base_path")),
	}

	if len(config.AwsKey) == 0 {
		return nil, fmt.Errorf("aws_key must be set")
	}

	if len(config.AwsSecret) == 0 {
		return nil, fmt.Errorf("aws_secret must be set")
	}

	if len(config.AwsRegion) == 0 {
		return nil, fmt.Errorf("aws_region must be set")
	}

	if len(config.GdprQueueUrl) == 0 {
		return nil, fmt.Errorf("gdpr_queue_url must be set")
	}
	if len(config.GdprCallbackQueueUrl) == 0 {
		return nil, fmt.Errorf("gdpr_callback_queue_url must be set")
	}
	if len(config.GdprBucket) == 0 {
		return nil, fmt.Errorf("gdpr_bucket must be set")
	}
	if len(config.GdprBasePath) == 0 {
		return nil, fmt.Errorf("gdpr_base_path must be set")
	}
	return config, nil
}

