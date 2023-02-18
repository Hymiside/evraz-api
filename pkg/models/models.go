package models

type ConfigServer struct {
	Host string
	Port string
}

type ConfigRepository struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}
