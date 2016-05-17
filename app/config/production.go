package config

func getProductionConfiguration() configKeyValues {

	return configKeyValues{

		BASE_WEB_URL:     "https://playment.in",
		BASE_API_URL:     "https://api.playment.in/",
		GOOGLE_CLIENT_ID: "some randomm id",
		SENDGRID_API_KEY: "s",
		PAYTM_AES:        "pay",
	}
}
