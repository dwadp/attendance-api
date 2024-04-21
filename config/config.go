package config

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Auth     Auth     `yaml:"auth"`
}

type Server struct {
	Timezone string `yaml:"timezone"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Auth struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type Database struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}
