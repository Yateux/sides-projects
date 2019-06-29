package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func TrackEvent(eventName string, params map[string]interface{}) error {
	// Data to send to analytics ms
	ret := map[string]interface{}{
		"params":    params,
		"client":    "TOSET",
		"eventName": eventName,
	}
	jsonValue, err := json.Marshal(ret)
	if err != nil {
		return errors.New(fmt.Sprintf("Analytics - Marshal failed: %s", err))
	}
	//HTTP call to analytics ms
	analyticsURL := os.Getenv("ANALYTICS_URL") + os.Getenv("DOMAIN")
	resp, err := http.Post(analyticsURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.New(fmt.Sprintf("Analytics - HTTP call failed:%s", err))
	}
	defer resp.Body.Close()
	return nil
}
