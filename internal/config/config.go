package config

type Config struct {
	db_url    string
	user_name string
}

func Read() (Config, error) {
	return Config{}, nil
}

func (cfg Config) SetUser() error {
	return nil
}
