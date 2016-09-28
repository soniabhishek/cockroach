package config

import (
	"errors"
	"github.com/crowdflux/angel/config_temp"
)

type configKey string

const (
	BASE_WEB_URL     = configKey("base.web")
	BASE_API_URL     = configKey("base.api")
	SENDGRID_API_KEY = configKey("sendgrid.apiKey")

	PAYTM_MID               = configKey("paytm.mid")
	PAYTM_SALES_WALLET_GUID = configKey("paytm.sales_wallet_guid")
	PAYTM_AES               = configKey("paytm.aes")

	FACEBOOK_APP_ID      = configKey("social.facebook.appId")
	FACEBOOK_APP_SECRET  = configKey("social.facebook.appSecret")
	GOOGLE_CLIENT_ID     = configKey("social.google.clientId")
	GOOGLE_CLIENT_SECRET = configKey("social.google.clientSecret")

	PG_DATABASE_NAME = configKey("postgres.database")
	PG_HOST          = configKey("postgres.host")
	PG_USERNAME      = configKey("postgres.username")
	PG_PASSWORD      = configKey("postgres.password")
	IP_ADDRESS       = configKey("ipAddress")

	AWS_REGION     = configKey("aws.region")
	AWS_SECRET_KEY = configKey("aws.secretKey")
	AWS_ACCESS_ID  = configKey("aws.accessId")
	S3_BUCKET      = configKey("aws.s3Bucket")

	MONGO_HOST    = configKey("mongo.host")
	MONGO_DB_NAME = configKey("mongo.name")

	CROWDY_BASE_API        = configKey("crowdy.base.api")
	CROWDY_AUTH_KEY        = configKey("crowdy.authKey")
	AUTHER_PLAYMENT_SECRET = configKey("auther.playment_secret")

	MONITOR_TIME_PERIOD         = configKey("fluMonitor.monitor_time_period")
	RETRY_TIME_PERIOD           = configKey("fluMonitor.retry_time_period")
	DEFAULT_FLU_THRESHOLD_COUNT = configKey("fluMonitor.default_flu_threshold_count")
	FLU_THRESHOLD_DURATION      = configKey("fluMonitor.flu_threshold_duration")
	FLU_RETRY_THRESHOLD         = configKey("fluMonitor.flu_retry_threshold")

	DOWNLOAD_PATH = configKey("workflow_step.manual_download_path")
	MEGATRON_API  = configKey("megatron.base.api")

	PLOG_LEVEL = configKey("plog.log_level")

	NEW_RELIC_KEY = configKey("newrelic.key")
)

// Gets the value for given key from the config file.
// It panics no configuration value is present
func (c configKey) Get() string {
	val := config_temp.Vipe.GetString(string(c))
	if val == "" {
		panic(errors.New("Configuration value not found [" + string(c) + "]"))
	}
	return val
}
