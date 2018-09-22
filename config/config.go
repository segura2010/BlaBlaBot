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
    Locale string // https://dev.blablacar.com/docs/versions/1.0/locales
    Currency string

    RefreshTime int // in minutes (every X minutes it will check the tasks and send the alerts)
    MaxTaskTime int // in hours (max time I can subscribe to a trip+date; example: If MaxTaskTime=10, I cannot subscribe to alert for trips that happens in more than 10 hours.)

    Whitelist []string // whitelist of telegram ids to allow use the bot (use "*" to allow everyone)
    
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
            Locale: "es_ES",
            Currency: "EUR",
            RefreshTime: 10, // 10 minutes
            MaxTaskTime: 240, // 240 hours (10 days),
            Whitelist: []string{"*"},
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

