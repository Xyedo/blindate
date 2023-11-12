package infrastructure

import (
	"flag"
	"os"
	"strconv"

	"github.com/invopop/validation"
	"github.com/joho/godotenv"
)

var Config config

type config struct {
	Host string
	Port int
	Env  string
	AWS  struct {
		AccessKeyId     string
		SecretAccessKey string
		BucketName      string
	}
	Postgre struct {
		Host     string
		Port     uint64
		User     string
		Password string
		Database string
	}
	Clerk struct {
		Token  string
		ApiKey string
		TestId string
	}
}

func (c config) validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.Env, validation.Required),

		validation.Field(&c.AWS.AccessKeyId, validation.Required),
		validation.Field(&c.AWS.SecretAccessKey, validation.Required),
		validation.Field(&c.AWS.BucketName, validation.Required),

		validation.Field(&c.Postgre.Host, validation.Required),
		validation.Field(&c.Postgre.Port, validation.Required),
		validation.Field(&c.Postgre.User, validation.Required),
		validation.Field(&c.Postgre.Password, validation.Required),
		validation.Field(&c.Postgre.Database, validation.Required),

		validation.Field(&c.Clerk.Token, validation.Required),
		validation.Field(&c.Clerk.ApiKey, validation.Required),
		validation.Field(&c.Clerk.TestId, validation.Required),
	)
}
func LoadConfig(filenames ...string) {
	_ = godotenv.Load(filenames...)

	// APP
	appPort, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	flag.StringVar(&Config.Host, "host", os.Getenv("APP_HOST"), "application host")
	flag.IntVar(&Config.Port, "port", appPort, "API server port")
	flag.StringVar(&Config.Env, "env", os.Getenv("APP_ENV"), "get env")

	//AWS s3
	flag.StringVar(&Config.AWS.AccessKeyId, "aws-access-key-id", os.Getenv("AWS_ACCESS_KEY_ID"), "S3 bucket name")
	flag.StringVar(&Config.AWS.SecretAccessKey, "aws-secret-access-key", os.Getenv("AWS_SECRET_ACCESS_KEY"), "S3 bucket name")
	flag.StringVar(&Config.AWS.BucketName, "s3-bucket-name", os.Getenv("AWS_BUCKET_NAME"), "S3 bucket name")

	// Postgre
	flag.StringVar(&Config.Postgre.Host, "db-host", os.Getenv("PG_HOST"), "PostgreSQL Host")
	var dbPort string
	flag.StringVar(&dbPort, "db-port", os.Getenv("PG_PORT"), "PostgreSQL Port")
	if v, err := strconv.ParseUint(dbPort, 10, 64); err == nil {
		Config.Postgre.Port = v
	}

	flag.StringVar(&Config.Postgre.User, "db-user", os.Getenv("PG_USER"), "PostgreSQL Username")
	flag.StringVar(&Config.Postgre.Password, "db-password", os.Getenv("PG_PASSWORD"), "PostgreSQL Password")
	flag.StringVar(&Config.Postgre.Database, "db-name", os.Getenv("PG_DB"), "PostgreSQL Database name")

	flag.StringVar(&Config.Clerk.Token, "clrek-auth-token", os.Getenv("CLREK_TOKEN"), "Clrek Token")
	flag.StringVar(&Config.Clerk.ApiKey, "clrek-apikey", os.Getenv("CLREK_API_KEY"), "Clrek ApiKey")
	flag.StringVar(&Config.Clerk.TestId, "clrek-test-id", os.Getenv("CLREK_TEST_ID"), "Clrek TestId")
	flag.Parse()

	err := Config.validate()
	if err != nil {
		panic(err)
	}

}
