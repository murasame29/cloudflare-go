package config

var Config *config

type config struct {
	CfAccountID        string `mapstructure:"CF_ACCOUNT_ID"`
	CfAccountKeyID     string `mapstructure:"CF_ACCOUNT_KEY_ID"`
	CfAccountKeySecret string `mapstructure:"CF_ACCOUNT_KEY_SECRET"`
	CfBucketName       string `mapstructure:"CF_BUCKET_NAME"`
	CfS3Endpoint       string `mapstructure:"CF_S3_ENDPOINT"`
}
