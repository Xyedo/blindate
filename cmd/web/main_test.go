package main_test

// var testQuery *sqlx.DB

// func TestMain(m *testing.M) {
// 	err := godotenv.Load("../../.env.dev")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	conn, err := sqlx.Open("postgres", os.Getenv("POSTGRE_DB_DSN_TEST"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	testQuery = conn
// 	gin.SetMode(gin.TestMode)
// 	os.Exit(m.Run())
// }

// func setupTestRouter(t *testing.T) http.Handler {
// 	var testCfg infra.Config
// 	flag.IntVar(&testCfg.Port, "port", 8080, "API server port")
// 	flag.StringVar(&testCfg.Env, "env", "development", "Environtment (development | staging | production)")

// 	flag.StringVar(&testCfg.DbConf.Dsn, "db-dsn", os.Getenv("POSTGRE_DB_DSN"), "PgSQL dsn")
// 	flag.IntVar(&testCfg.DbConf.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
// 	flag.IntVar(&testCfg.DbConf.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
// 	flag.StringVar(&testCfg.DbConf.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

// 	flag.StringVar(&testCfg.Token.AccessSecret, "jwt-access-secret", os.Getenv("JWT_ACCESS_SECRET_KEY"), "Jwt Access")
// 	flag.StringVar(&testCfg.Token.RefreshSecret, "jwt-refresh-secret", os.Getenv("JWT_REFRESH_SECRET_KEY"), "Jwt Access")
// 	flag.StringVar(&testCfg.Token.AccessExpires, "jwt-access-expires", os.Getenv("JWT_ACCESS_EXPIRES"), "Jwt Access")
// 	flag.StringVar(&testCfg.Token.RefreshExpires, "jwt-refresh-expires", os.Getenv("JWT_REFRESH_EXPIRES"), "Jwt Access")

// 	route := testCfg.Container(testQuery)
// 	return api.Routes(route)
// }
// func cleanUp() {
// 	testQuery.MustExec(`DELETE FROM users WHERE 1=1`)
// 	testQuery.MustExec(`DELETE FROM authentications WHERE 1=1`)
// }
