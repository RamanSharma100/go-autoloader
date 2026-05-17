package core

type Config struct {
	AutoParse bool
}

func DefaultConfig() Config {
	return Config{
		AutoParse: true,
	}
}
