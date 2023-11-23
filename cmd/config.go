package main

type serviceConfig struct {
	DBConfig dbConfig `yaml:"database"`
}

type dbConfig struct {
	Port     int    `yaml:"port"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	StrConn  string `yaml:"strConn"`
	DBName   string `yaml:"dbName"`
}
