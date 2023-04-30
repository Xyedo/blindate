package infrastructure

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host       string
	Port       int
	Env        string
	BucketName string
	DbConf     struct {
		Host         string
		Port         string
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
}

func (cfg *Config) LoadConfig(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Panic(err)
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Panic(err)
	}

	flag.StringVar(&cfg.Host, "host", os.Getenv("APP_HOST"), "application host")
	flag.IntVar(&cfg.Port, "port", appPort, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environtment (development | staging | production)")

	flag.StringVar(&cfg.DbConf.Host, "db-host", os.Getenv("PG_HOST"), "PostgreSQL Host")
	flag.StringVar(&cfg.DbConf.Port, "db-port", os.Getenv("PG_PORT"), "PostgreSQL Port")
	flag.StringVar(&cfg.DbConf.User, "db-user", os.Getenv("PG_USER"), "PostgreSQL Username")
	flag.StringVar(&cfg.DbConf.Password, "db-password", os.Getenv("PG_PASSWORD"), "PostgreSQL Password")
	flag.StringVar(&cfg.DbConf.Database, "db-name", os.Getenv("PG_DB"), "PostgreSQL Database name")

	flag.IntVar(&cfg.DbConf.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DbConf.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DbConf.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.Token.AccessSecret, "jwt-access-secret", os.Getenv("JWT_ACCESS_SECRET_KEY"), "Jwt Access")
	flag.StringVar(&cfg.Token.RefreshSecret, "jwt-refresh-secret", os.Getenv("JWT_REFRESH_SECRET_KEY"), "Jwt Access")
	flag.StringVar(&cfg.Token.AccessExpires, "jwt-access-expires", os.Getenv("JWT_ACCESS_EXPIRES"), "Jwt Access")
	flag.StringVar(&cfg.Token.RefreshExpires, "jwt-refresh-expires", os.Getenv("JWT_REFRESH_EXPIRES"), "Jwt Access")

	flag.StringVar(&cfg.BucketName, "s3-bucket-name", os.Getenv("AWS_BUCKET_NAME"), "S3 bucket name")
	flag.Parse()
}
