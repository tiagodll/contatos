package util

import "os"

type Config struct {
	AppPort            string
	AppUrl             string
	DbConnectionString string
	LogPath            string
	StaticPath         string
	PnorcoID           string
	Jwt                JwtConfig
}

type JwtConfig struct {
	Secret   string
	Issuer   string
	Audience string
}

func GetEnvConfig() Config {
	config := Config{
		AppPort:            os.Getenv("APP_PORT"),
		AppUrl:             os.Getenv("APP_URL"),
		DbConnectionString: os.Getenv("CONNECTION_STRING"),
		LogPath:            os.Getenv("LOG_PATH"),
		StaticPath:         os.Getenv("STATIC_PATH"),
		PnorcoID:           os.Getenv("PNORCO_ID"),
		Jwt: JwtConfig{
			Secret:   os.Getenv("JWT_SECRET"),
			Issuer:   os.Getenv("JWT_ISSUER"),
			Audience: os.Getenv("JWT_AUDIENCE"),
		},
	}
	return config
}

func (c Config) GetLoginUrl() string {
	return c.Jwt.Issuer + "/login?app_id=" + c.PnorcoID
}
