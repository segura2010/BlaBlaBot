package bot


import (
    "time"
    "fmt"
    "strings"
    //"log"
    // "strconv"
    // "encoding/hex"
    
    "gopkg.in/tucnak/telebot.v1"

    "BlaBlaBot/redisdb"
    "BlaBlaBot/config"
)

type TelegramBot struct{
	Bot *telebot.Bot
    Token string
    Started bool
}

var instance *TelegramBot = nil

func CreateInstance(token string) *TelegramBot {
    instance = &TelegramBot{Token:token, Started:false}
    bot, err := telebot.NewBot(token)
    if err != nil {
        panic(err)
    }

    instance.Started = true
    instance.Bot = bot
    go listenMessages()

    return instance
}

func GetInstance() *TelegramBot {
    return instance
}

func RefreshSession(){
    CreateInstance(instance.Token)
}

func SendMessage(to int64, message string, options *telebot.SendOptions){
    myBot := GetInstance()

    chat := telebot.Chat{ID: to}
    myBot.Bot.SendMessage(chat, message, options)
}

func listenMessages(){
    myBot := GetInstance()
    messages := make(chan telebot.Message)
    myBot.Bot.Listen(messages, 1*time.Second)

    sendOptions := telebot.SendOptions{
        ParseMode: "Markdown",
    }

    cfg := config.GetInstance()

    for message := range messages {
        // log.Printf("Received message..")
        userId := fmt.Sprintf("%d", message.Sender.ID)
        // chatId := fmt.Sprintf("%d", message.Chat.ID)

        if cfg.Maintenance.Enabled {
            myBot.Bot.SendMessage(message.Chat, cfg.Maintenance.Description, nil)
        }else if strings.Index(message.Text, "/add ") == 0 {
            task := message.Text[5:]
            formatedTask := strings.Split(task, " ")
            if len(formatedTask) != 3{
                myBot.Bot.SendMessage(message.Chat, "Bad format :(", nil)
                continue
            }

            if isValidDate(formatedTask[0]){
                err := redisdb.AddTask(userId, formatedTask[0], formatedTask[1], formatedTask[2])
                if err != nil{
                    myBot.Bot.SendMessage(message.Chat, "There was an error adding your subscription :(", nil)
                    continue
                }
                myBot.Bot.SendMessage(message.Chat, "Your subscription was added!", nil)
            }else{
                myBot.Bot.SendMessage(message.Chat, "Invalid date", nil)
            }
        }else if strings.Index(message.Text, "/delete ") == 0 {
            date := message.Text[8:]
            err := redisdb.DeleteTask(userId, date)
            if err != nil{
                myBot.Bot.SendMessage(message.Chat, "There was an error deleting your subscription :(", nil)
                continue
            }
            myBot.Bot.SendMessage(message.Chat, "Your subscription was deleted!", nil)
        }else if strings.Index(message.Text, "/me") == 0 {
            tasks, err := redisdb.GetUserTasks(userId)
            if err != nil{
                myBot.Bot.SendMessage(message.Chat, "There was an error retrieving your subscriptions :(", nil)
                continue
            }
            SendMessage(message.Chat.ID, formatTasks(tasks), &sendOptions)
        }else{
            // help..
            r := fmt.Sprintf("Available commands:")
            r += fmt.Sprintf("\n\t/add date departure arrival: subscrite to alerts for trips in an specific date (format: YYYY-MM-DD)")
            r += fmt.Sprintf("\n\t/delete date: delete your subscription to alert for an specific date")
            r += fmt.Sprintf("\n\t/me: get your subscriptions")
            myBot.Bot.SendMessage(message.Chat, r, nil)
        }
    }
}

func formatTasks(taskId []string)(string){
    result := fmt.Sprintf("*Subscriptions*:\n")

    for _, t := range taskId{
        task, _ := redisdb.GetTaskByKey(t)
        splittedTask := strings.Split(t, ":")
        date := splittedTask[2]
        places := strings.Split(task, ":")
        result += fmt.Sprintf("*%s*: %s - %s", date, places[0], places[1])
    }

    return result
}

func FormatTrip(trip blablacarapi.Trip)(string){
    result := fmt.Sprintf("*New trip*:\n")

    result += fmt.Sprintf("*%s*", trip.Departure_Date)
    // result += fmt.Sprintf("*From:*", trip.Departure_Place)

    return result
}

func isValidDate(date string) (bool){
    return true
}


