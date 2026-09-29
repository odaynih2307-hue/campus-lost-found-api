package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port               string
	Env                string
	DatabaseURL        string
	JWTSecret          string
	JWTIssuer          string
	AccessTokenMinutes int
	RefreshTokenDays   int
	DBMaxConns         int32
	DBMinConns         int32
}

func Load() Config {
	return Config{
		Port: getenv("APP_PORT", "3000"), Env: getenv("APP_ENV", "development"),
		DatabaseURL:        getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/campus_lost_found"),
		JWTSecret:          getenv("JWT_SECRET", "dev-secret-change-this-to-a-long-random-secret"),
		JWTIssuer:          getenv("JWT_ISSUER", "campus-lost-found-api"),
		AccessTokenMinutes: getenvInt("ACCESS_TOKEN_MINUTES", 15), RefreshTokenDays: getenvInt("REFRESH_TOKEN_DAYS", 7),
		DBMaxConns: int32(getenvInt("DB_MAX_CONNS", 10)), DBMinConns: int32(getenvInt("DB_MIN_CONNS", 2)),
	}
}
func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func getenvInt(k string, d int) int {
	v, err := strconv.Atoi(os.Getenv(k))
	if err != nil || v <= 0 {
		return d
	}
	return v
}
