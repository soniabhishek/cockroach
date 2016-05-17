package config

func getDevelopmentConfiguration() configKeyValues {

	return configKeyValues{

		BASE_WEB_URL:     "localhost/",
		BASE_API_URL:     "localhost:8999/",
		GOOGLE_CLIENT_ID: "some randomm id",
		PAYTM_AES:        "pay",
		SENDGRID_API_KEY: "s",
		AWS_REGION:       "ap-southeast-1",
		AWS_SECRET_KEY:   "v0/dr/aXopQKwfwVdDhhZJlNVes/2VzjYQCnfeKS",
		AWS_ACCESS_ID:    "AKIAILPB5FXULHOVAO3Q",
		S3_BUCKET:        "playmentdevelopment",
		MONGO_HOST:       "localhost:27017",
		MONGO_DB_NAME:    "playment_mongo_local",
		DB_DATABASE_NAME: "playment_local",
		DB_HOST:          "localhost:5432",
	}
}
