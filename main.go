package main

import (
    "log"
    "fmt"
    "flag"
    "time"
    "strconv"

    "BlaBlaBot/bot"
    "BlaBlaBot/blablacarapi"
    "BlaBlaBot/redisdb"
    "BlaBlaBot/config"
)

/* Update version number on each release:
    Given a version number x.y.z, increment the:

    x - major release
    y - minor release
    z - build number
*/
const GOSTOCK_VERSION = "0.0.1"
var Commit string
var CompilationDate string

func showVersionInfo(){
    fmt.Println("----------------------------------------")
    fmt.Printf("BlaBlaBot v%s\nCommit: %s\nCompilation date: %s\n", GOSTOCK_VERSION, Commit, CompilationDate)
    fmt.Println("----------------------------------------")
}

func main(){
    // Command line options
    start := flag.Bool("s", false, "Start bot")
    cfgFile := flag.String("c", "./config.json", "Config file")
    version := flag.Bool("v", false, "Show current BlaBlaBot version")
    flag.Parse()

    log.Printf("\nLoading config..")
    config.CreateInstance(*cfgFile)
    
    if *version{
        showVersionInfo()
    }

    if *start{
        log.Printf("Started!")
        instance := bot.CreateInstance(config.GetInstance().TelegramBot.Token)
        if instance == nil{
            panic("Unable to create TelegramBot")
        }
        
        sendOptions := telebot.SendOptions{
            ParseMode: "Markdown",
        }
        
        // check for trips to send alerts :)
        for{
            // listen for messages...
            tasks, _ := redisdb.GetTasks()
            for _, t := range tasks{
                task, _ := redisdb.GetTaskByKey(t)
                splittedTask := strings.Split(t, ":")
                date := splittedTask[2]
                places := strings.Split(task, ":")

                trips, err := blablacarapi.GetTrips(places[0], places[1], "es_ES", "EUR", date)
                if err != nil{
                    log.Printf("Error getting trips: %s", err)
                    continue
                }

                // send alert for trip
                uId, err := strconv.Atoi(splittedTask[1])
                if err != nil{
                    continue
                }
                bot.SendMessage(uId, , &sendOptions)

                // save the sent alert to avoid sending it in the future !!
            }
            time.Sleep(time.Duration(5) * time.Minute)
        }
    }
    
}

