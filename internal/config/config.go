package config

var (
	cfg Config
)

func init() {
	v, err := load()
	if err != nil {
		panic(err)
	}
	if err = v.Unmarshal(&cfg); err != nil {
		panic(err)
	}
}

func Get() *Config {
	return &cfg
}
