package svr

type Config struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func NewDefaultConfig() Config {
	return Config{
		Address: "localhost",
		Port:    1000,
	}
}
