package config

type Config struct {
    Host string
    Port string
    User string
    Password string
    Database string
}

//User-customizable settings
var Settings = Config{
    "0.0.0.0",
    "8080",
    "db_user",
    "password",
    "db",
}