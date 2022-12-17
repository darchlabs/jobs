package config

type Config struct {
	DatabaseURL string `envconfig:"database_filepath" required:"true"`
	Port        string `envconfig:"port" required:"true"`
	PrivateKey  string `envconfig:"private_key" required:"true"`
}
