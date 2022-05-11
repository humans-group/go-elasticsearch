package es

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Error struct {
	StatusCode int
	Type       string
	Reason     string
}

func (e Error) Error() string {
	return fmt.Sprintf("Error: [%d] %s: %s",
		e.StatusCode,
		e.Type,
		e.Reason,
	)
}

func ExtractError(res *esapi.Response) error {
	if res.IsError() {
		var raw map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
			return fmt.Errorf("Failure to parse response body: %w", err)
		} else {

			typ, ok := raw["error"].(map[string]interface{})["type"].(string)
			if !ok {
				return fmt.Errorf("Failure to parse type to string")
			}

			reason, ok := raw["error"].(map[string]interface{})["reason"].(string)
			if !ok {
				return fmt.Errorf("Failure to parse reason to string")
			}

			return Error{
				StatusCode: res.StatusCode,
				Type:       typ,
				Reason:     reason,
			}
		}
	}
	return nil
}
