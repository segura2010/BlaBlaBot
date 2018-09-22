package redisdb

import (

    "github.com/go-redis/redis"

    "BlaBlaBot/config"
)

var instance *redis.Client = nil

func GetInstance() *redis.Client {
    if instance == nil {
    	host := config.GetInstance().RedisDB.Host
    	port := config.GetInstance().RedisDB.Port
    	addr := host + ":" + port
        instance = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: config.GetInstance().RedisDB.Pass,
			DB:       0,  // use default DB
		})

		_, err := instance.Ping().Result()
		if err != nil {
			panic(err)
		}
    }

    return instance
}

// A task is an user trip alert for an specific date; for example:
// task:USERID:YYYY-MM-DD
// and the value of the task is the departure and arrival place; example:
// Baza:Málaga
// Before check new trips for the task, we check if the date has passed, if it passed, I have to delete ot from the DB.
func TaskPreffix() string{
	return "task:"
}

// An alert element is used to remember if we sent an alert for a trip or not
// The structure is:
// alert:USERID:YYYY-MM-DD:TRIPID
func AlertPreffix() string{
	return "alert:"
}


func AddTask(userid, date, departure, arrival string) (error){
	err := GetInstance().Set(TaskPreffix() + userid +":"+ date, departure +":"+ arrival, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(userid, date string) (error){
	err := GetInstance().Del(TaskPreffix() + userid +":"+ date).Err()
	if err != nil {
		return err
	}

	return nil
}
func DeleteTaskByKey(key string) (error){
	err := GetInstance().Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}


func GetTask(userid, date string) (string, error){
	task, err := GetInstance().Get(TaskPreffix() + userid +":"+ date).Result()
	return task, err
}
func GetTaskByKey(key string) (string, error){
	task, err := GetInstance().Get(key).Result()
	return task, err
}

func GetUserTasks(userid string) ([]string, error){
	tasks, err := GetInstance().Keys(TaskPreffix() + userid +":*").Result()
	return tasks, err
}

func GetTasks() ([]string, error){
	tasks, err := GetInstance().Keys(TaskPreffix() + "*").Result()
	return tasks, err
}

func AddAlert(key string) error{
	err := GetInstance().Set(AlertPreffix() + key, "sent", 0).Err()
	if err != nil {
		return err
	}

	return nil
}
func GetAlertByKey(key string) (string, error){
	alert, err := GetInstance().Get(key).Result()
	return alert, err
}


func DeleteKey(key string) (error){
	err := GetInstance().Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}

// ONLY CALL THIS FUNCTION WHEN TASK_DATE < NOW
// It remove everything related to a task in the DB
// Because we do not need to save the data once the date has passed
func DeleteAllTaskRelatedStuff(uid, date string){
	tasks, _ := GetInstance().Keys("*:"+ uid +":"+ date).Result()
	for _, t := range tasks{
		DeleteKey(t)
	}

	alerts, _ := GetInstance().Keys("*:"+ uid +":"+ date +":*").Result()
	for _, a := range alerts{
		DeleteKey(a)
	}
}

func Flush() (error){
	return GetInstance().FlushDB().Err()
}


