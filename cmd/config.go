package main

type serviceConfig struct {
	DBConfig    dbConfig    `yaml:"database"`
	RedisConfig redisConfig `yaml:"redis"`
}

type dbConfig struct {
	Port     int    `yaml:"port"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	StrConn  string `yaml:"strConn"`
	DBName   string `yaml:"dbName"`
}
type redisConfig struct {
	Port     int    `yaml:"port"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DBName   int    `yaml:"dbName"`
	StrConn  string `yaml:"strConn"`
}
