package config

func getStagingConfiguration() configKeyValues {

	return configKeyValues{

		BASE_WEB_URL:     "staging.playment.in/",
		BASE_API_URL:     "stagingapi.playment.in/",
		AWS_REGION:       "ap-southeast-1",
		AWS_SECRET_KEY:   "v0/dr/aXopQKwfwVdDhhZJlNVes/2VzjYQCnfeKS",
		AWS_ACCESS_ID:    "AKIAILPB5FXULHOVAO3Q",
		S3_BUCKET:        "playmentproduction",
		MONGO_HOST:       "localhost:27017",
		MONGO_DB_NAME:    "playment_mongo_staging",
		DB_DATABASE_NAME: "playment_staging",
		DB_HOST:          "localhost:5432",
	}
}
