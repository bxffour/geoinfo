package main

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	port int
	env  string
	db   struct {
		user         string
		port         int
		password     string
		dbname       string
		host         string
		sslrootcert  string
		sslcert      string
		sslkey       string
		sslmode      string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

func setupFlags() error {
	flag.String("config", "", "path to config file")

	flag.Int("port", 8080, "server port")
	flag.String("env", "development", "Environment ( development | staging | production )")

	// Postgres connection parameters
	flag.String("db.user", "", "Postgresql database user")
	flag.String("db.password", "", "Postgresql database password")
	flag.String("db.dbname", "crest-countries", "Postgresql database name")
	flag.Int("db.port", 5432, "Postgresql port")
	flag.String("db.host", "localhost", "Postgresql host")
	flag.String("db.sslcert", "", "Postgresql client certificate file")
	flag.String("db.sslkey", "", "Postgresql client key file")
	flag.String("db.sslrootcert", "", "Postgresql ssl certificate authority")
	flag.String("db.sslmode", "prefer", "Postgresql sslmode (disable | allow | prefer | require | verify-ca | verify-full)")

	flag.Int("db.max-open-conns", 25, "Postgresql maximum open connections")
	flag.Int("db.max-idle-conns", 25, "Postgresql maximum idle connections")
	flag.String("db.max-idle-time", "15m", "Postgresql maximum idle time")

	flag.Float64("limiter.rps", 2, "Rate limiter maximum request per second")
	flag.Int("limiter.burst", 4, "Rate limiter maximum burst")
	flag.Bool("limiter.enabled", false, "Enable rate limiter")

	flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	return viper.BindPFlags(flag.CommandLine)
}

func setupConfig(c *config) error {
	var err error

	config := viper.GetString("config")
	viper.SetConfigFile(config)

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	c.port = viper.GetInt("port")
	c.env = viper.GetString("env")

	c.db.user = viper.GetString("db.user")
	c.db.password = viper.GetString("db.password")
	c.db.dbname = viper.GetString("db.dbname")
	c.db.port = viper.GetInt("db.port")
	c.db.host = viper.GetString("db.host")
	c.db.sslrootcert = viper.GetString("db.sslrootcert")
	c.db.sslcert = viper.GetString("db.sslcert")
	c.db.sslkey = viper.GetString("db.sslkey")
	c.db.sslmode = viper.GetString("db.sslmode")

	c.db.maxOpenConns = viper.GetInt("db.max-open-conns")
	c.db.maxIdleConns = viper.GetInt("db.max-idle-conns")
	c.db.maxIdleTime = viper.GetString("db.max-idle-time")

	c.limiter.rps = viper.GetFloat64("limiter.rps")
	c.limiter.burst = viper.GetInt("limiter.burst")
	c.limiter.enabled = viper.GetBool("limiter.enabled")

	return nil
}