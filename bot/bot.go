package bot


import (
    "time"
    "fmt"
    "strings"
    //"log"
    // "strconv"
    // "encoding/hex"
    
    "gopkg.in/tucnak/telebot.v1"

    "BlaBlaBot/blablacarapi"
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

        if !isInTheWhitelist(userId){
            myBot.Bot.SendMessage(message.Chat, "This is a *private* bot and you are not in the *whitelist* 😢\nContact with the administrator to use me.\n\nYour Telegram ID: "+ userId, &sendOptions)
            continue
        }

        if cfg.Maintenance.Enabled {
            myBot.Bot.SendMessage(message.Chat, cfg.Maintenance.Description, nil)
        }else if strings.Index(message.Text, "/add ") == 0 {
            task := message.Text[5:]
            formatedTask := strings.Split(task, " ")
            if len(formatedTask) != 3{
                myBot.Bot.SendMessage(message.Chat, "Bad format :(", nil)
                continue
            }

            if isValidDate(formatedTask[2]){
                err := redisdb.AddTask(userId, formatedTask[2], strings.ToUpper(formatedTask[0]), strings.ToUpper(formatedTask[1]))
                if err != nil{
                    myBot.Bot.SendMessage(message.Chat, "There was an error adding your subscription 😢", nil)
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
                myBot.Bot.SendMessage(message.Chat, "There was an error deleting your subscription 😢", nil)
                continue
            }
            myBot.Bot.SendMessage(message.Chat, "Your subscription was deleted!", nil)
        }else if strings.Index(message.Text, "/me") == 0 {
            tasks, err := redisdb.GetUserTasks(userId)
            if err != nil{
                myBot.Bot.SendMessage(message.Chat, "There was an error retrieving your subscriptions 😢", nil)
                continue
            }
            SendMessage(message.Chat.ID, formatTasks(tasks), &sendOptions)
        }else{
            // help..
            r := fmt.Sprintf("Available commands:")
            r += fmt.Sprintf("\n\t/add departure arrival date: subscrite to alerts for trips in an specific date (format: YYYY-MM-DD)")
            r += fmt.Sprintf("\n\t/delete date: delete your subscription to alert for an specific date")
            r += fmt.Sprintf("\n\t/me: get your subscriptions")
            myBot.Bot.SendMessage(message.Chat, r, nil)
        }
    }
}

func SendTripAlert(uId int64, trip blablacarapi.Trip){
    sendOptions := telebot.SendOptions{
        ParseMode: "Markdown",
    }

    SendMessage(uId, formatTrip(trip), &sendOptions)
}

func formatTasks(taskId []string)(string){
    result := fmt.Sprintf("🚙 *Subscriptions*\n")

    for _, t := range taskId{
        task, _ := redisdb.GetTaskByKey(t)
        splittedTask := strings.Split(t, ":")
        date := splittedTask[2]
        places := strings.Split(task, ":")
        result += fmt.Sprintf("\n📅*%s*: %s ➡️ %s", date, places[0], places[1])
    }

    return result
}

func formatTrip(trip blablacarapi.Trip)(string){
    result := fmt.Sprintf("😊*New trip*:\n")

    result += fmt.Sprintf("\n📅*%s*", trip.Departure_Date)
    result += fmt.Sprintf("\n*From:* %s (%s)", trip.Departure_Place.Address, trip.Departure_Place.City_Name)
    result += fmt.Sprintf("\n*To:* %s (%s)", trip.Arrival_Place.Address, trip.Arrival_Place.City_Name)
    result += fmt.Sprintf("\n*Price:* %s", trip.Price_With_Commission.String_Value)
    result += fmt.Sprintf("\n\n🔗%s", trip.Links["_front"])

    return result
}

func isValidDate(date string) (bool){
    now := time.Now()
    taskDate, err := time.Parse("2006-01-02", date)
    if err != nil{
        return false
    }

    diff := taskDate.Sub(now)
    return (diff > 0) && (diff.Hours() < float64(config.GetInstance().MaxTaskTime))
}

func isInTheWhitelist(uid string) (bool){
    whitelist := config.GetInstance().Whitelist
    for _, u := range whitelist{
        if u == uid || u == "*"{
            return true
        }
    }
    return false
}


