package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Verify user exists in user-service
func VerifyUser(userID int) (bool, error) {
	url := fmt.Sprintf("http://localhost:8081/api/v1/users/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("User not found: %d", userID)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return false, err
	}

	return true, nil
}
