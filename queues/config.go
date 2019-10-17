package queues

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	AwsKey               string
	AwsSecret            string
	AwsRegion            string
	GdprQueueUrl         string
	CmsQueueUrl          string
	GdprCallbackQueueUrl string
	GdprBucket           string
	GdprBasePath         string
	SqsMaxMessages       int
	TimeWaitSeconds      int
}

func InitConfig() (*Config, error) {
	config := &Config{
		AwsKey:               viper.GetString("aws_key"),
		AwsSecret:            viper.GetString("aws_secret"),
		AwsRegion:            viper.GetString("aws_region"),
		GdprQueueUrl:         viper.GetString("gdpr_queue_url"),
		CmsQueueUrl:          viper.GetString("cms_queue_url"),
		GdprCallbackQueueUrl: viper.GetString("gdpr_callback_queue_url"),
		GdprBucket:           viper.GetString("gdpr_bucket"),
		GdprBasePath:         viper.GetString("gdpr_base_path"),
		SqsMaxMessages:       viper.GetInt("sqs_max_messages"),
		TimeWaitSeconds:      viper.GetInt("sqs_max_messages"),
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

	if len(config.CmsQueueUrl) == 0 {
		return nil, fmt.Errorf("cms_queue_url must be set")
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

	if config.SqsMaxMessages == 0 {
		config.SqsMaxMessages = 1
	}

	if config.TimeWaitSeconds == 0 {
		config.TimeWaitSeconds = 1
	}
	return config, nil
}
