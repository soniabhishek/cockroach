package config

import (
	"errors"
)

type configKey string

const (
	BASE_WEB_URL            = configKey("base.web")
	BASE_API_URL            = configKey("base.api")
	SENDGRID_API_KEY        = configKey("sendgrid.apiKey")
	PAYTM_MID               = configKey("paytm.mid")
	PAYTM_SALES_WALLET_GUID = configKey("paytm.sales_wallet_guid")
	PAYTM_AES               = configKey("paytm.aes")
	FACEBOOK_APP_ID         = configKey("social.facebook.appId")
	FACEBOOK_APP_SECRET     = configKey("social.facebook.appSecret")
	GOOGLE_CLIENT_ID        = configKey("social.google.clientId")
	GOOGLE_CLIENT_SECRET    = configKey("social.google.clientSecret")
	PG_DATABASE_NAME        = configKey("postgres.database")
	PG_HOST                 = configKey("postgres.host")
	PG_USERNAME             = configKey("postgres.username")
	PG_PASSWORD             = configKey("postgres.password")
	IP_ADDRESS              = configKey("ipAddress")
	AWS_REGION              = configKey("aws.region")
	AWS_SECRET_KEY          = configKey("aws.secretKey")
	AWS_ACCESS_ID           = configKey("aws.accessId")
	S3_BUCKET               = configKey("aws.s3Bucket")
	MONGO_HOST              = configKey("mongo.host")
	MONGO_DB_NAME           = configKey("mongo.name")
	CROWDY_BASE_API         = configKey("crowdy.base.api")
	CROWDY_AUTH_KEY         = configKey("crowdy.authKey")
	MONITOR_TIME_PERIOD     = configKey("fluMonitor.monitor_time_period")
	RETRY_TIME_PERIOD       = configKey("fluMonitor.retry_time_period")
	FLU_THRESHOLD_COUNT     = configKey("fluMonitor.flu_threshold_count")
	FLU_THRESHOLD_DURATION  = configKey("fluMonitor.flu_threshold_duration")
	DOWNLOAD_PATH           = configKey("workflow_step.manual_download_path")
)

// Gets the value for given key from the config file.
// It accepts only configKey type, which is private, so
// only the above consts can be passed
// It panics no configuration value is present
func Get(c configKey) string {
	val := configProvider.GetString(string(c))
	if val == "" {
		panic(errors.New("Configuration value not found"))
	}
	return val
}
