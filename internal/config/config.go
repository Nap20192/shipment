package config

type Config struct {
	DBHost      string `env:"DB_HOST" default:"localhost"`
	DBPort      int    `env:"DB_PORT" default:"5432"`
	DBUser      string `env:"DB_USER" default:"postgres"`
	DBPassword  string `env:"DB_PASSWORD"`
	DBName      string `env:"DB_NAME" default:"shipment"`
	LogLevel    string `env:"LOG_LEVEL" default:"info"`
	LogDir      string `env:"LOG_DIR" default:"./logs"`
}
