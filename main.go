package main

import (
    "log"
    "fmt"
    "flag"
    "time"
    "strconv"
    "strings"

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
        
        // check for trips to send alerts :)
        for{
            // listen for messages...
            cache := make(map[string]blablacarapi.TripsResponse)
            tasks, _ := redisdb.GetTasks()
            for _, t := range tasks{
                task, _ := redisdb.GetTaskByKey(t)
                splittedTask := strings.Split(t, ":")
                date := splittedTask[2]
                places := strings.Split(task, ":")
                uId, err := strconv.Atoi(splittedTask[1])
                if err != nil{
                    continue
                }

                // check if trip passed
                now := time.Now()
                taskDate, err := time.Parse("2006-01-02", date)
                if err != nil{
                    continue
                }
                diff := now.Sub(taskDate)
                if diff > 0{
                    // task time passed, delete
                    log.Printf("Task %s date passed, deleting everything", t)
                    go redisdb.DeleteAllTaskRelatedStuff(date)
                    continue
                }

                var trips blablacarapi.TripsResponse
                exists := false

                if trips, exists = cache[date +":"+ task]; !exists{
                    // non cached task
                    trips, err = blablacarapi.GetTrips(places[0], places[1], config.GetInstance().Locale, config.GetInstance().Currency, date)
                    if err != nil{
                        log.Printf("Error getting trips: %s", err)
                        continue
                    }
                    // cache trip
                    cache[date +":"+ task] = trips
                }else{
                    log.Printf("Task %s already cached!", (date +":"+ task))
                }

                for _, trip := range trips.Trips{
                    // check if already sent
                    alreadySent, _ := redisdb.GetAlertByKey(redisdb.AlertPreffix() + splittedTask[1] +":"+ date +":"+ trip.Permanent_Id)
                    if alreadySent != ""{
                        log.Printf("Alert %s already sent", trip.Permanent_Id)
                        continue
                    }

                    // send alert for trip
                    bot.SendTripAlert(int64(uId), trip)

                    // save the sent alert to avoid sending it in the future !!
                    redisdb.AddAlert(splittedTask[1] +":"+ date +":"+ trip.Permanent_Id)
                }
            }
            time.Sleep(time.Duration(config.GetInstance().RefreshTime) * time.Minute)
        }
    }
    
}

