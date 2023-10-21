package infrastructure

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Config config

type config struct {
	Host       string
	Port       int
	Env        string
	BucketName string
	DbConf     struct {
		Host         string
		Port         uint64
		User         string
		Password     string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Token struct {
		AccessSecret   string
		RefreshSecret  string
		AccessExpires  string
		RefreshExpires string
	}
	ClrekToken string
}

func LoadConfig(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Panic(err)
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Panic(err)
	}

	flag.StringVar(&Config.Host, "host", os.Getenv("APP_HOST"), "application host")
	flag.IntVar(&Config.Port, "port", appPort, "API server port")
	flag.StringVar(&Config.Env, "env", "development", "Environtment (development | staging | production)")

	flag.StringVar(&Config.DbConf.Host, "db-host", os.Getenv("PG_HOST"), "PostgreSQL Host")
	var dbPort string
	flag.StringVar(&dbPort, "db-port", os.Getenv("PG_PORT"), "PostgreSQL Port")
	flag.StringVar(&Config.DbConf.User, "db-user", os.Getenv("PG_USER"), "PostgreSQL Username")
	flag.StringVar(&Config.DbConf.Password, "db-password", os.Getenv("PG_PASSWORD"), "PostgreSQL Password")
	flag.StringVar(&Config.DbConf.Database, "db-name", os.Getenv("PG_DB"), "PostgreSQL Database name")

	flag.IntVar(&Config.DbConf.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&Config.DbConf.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&Config.DbConf.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&Config.Token.AccessSecret, "jwt-access-secret", os.Getenv("JWT_ACCESS_SECRET_KEY"), "Jwt Access")
	flag.StringVar(&Config.Token.RefreshSecret, "jwt-refresh-secret", os.Getenv("JWT_REFRESH_SECRET_KEY"), "Jwt Access")
	flag.StringVar(&Config.Token.AccessExpires, "jwt-access-expires", os.Getenv("JWT_ACCESS_EXPIRES"), "Jwt Access")
	flag.StringVar(&Config.Token.RefreshExpires, "jwt-refresh-expires", os.Getenv("JWT_REFRESH_EXPIRES"), "Jwt Access")

	flag.StringVar(&Config.BucketName, "s3-bucket-name", os.Getenv("AWS_BUCKET_NAME"), "S3 bucket name")
	flag.StringVar(&Config.ClrekToken, "clrek-auth-token", os.Getenv("CLREK_TOKEN"), "Clrek Token")
	flag.Parse()

	if v, err := strconv.ParseUint(dbPort, 10, 64); err == nil {
		Config.DbConf.Port = v
	}
}
