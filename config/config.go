package config

import (
    "log"
    "encoding/json"
    "io/ioutil"
)

// MyConfig struct
// This is the struct that the config.json must have
type MyConfig struct {
    ApiKey string // BlaBlaCar API Key
    
    // DB info
    RedisDB RedisConfig

    // Telegram info
    TelegramBot TelegramConfig
    // Maintenance Info
    Maintenance MaintenanceConfig
}

// DB Config
type RedisConfig struct {
    Host string
    Port string
    Pass string
}

// Telegram Config
type TelegramConfig struct {
    Token string
}

// Maintenance Mode Config
type MaintenanceConfig struct {
    Enabled bool
    Description string
}

var instance *MyConfig = nil

func CreateInstance(filename string) *MyConfig {
    var err error
    instance, err = loadConfig(filename)
    log.Printf("Error loading config file: %s\nUsing default config.", err)
    if err != nil {
        // use defaults
        instance = &MyConfig{
            ApiKey: "",
            RedisDB: RedisConfig{
                Host: "localhost",
                Port: "6379",
                Pass: "",
            },
        }
    }

    return instance
}

func GetInstance() *MyConfig {
    return instance
}

func loadConfig(filename string) (*MyConfig, error){
    var s *MyConfig

    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return s, err
    }
    // Unmarshal json
    err = json.Unmarshal(bytes, &s)
    return s, err
}

