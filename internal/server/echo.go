package server

import (
	"encoding/json"
	"errors"
)

type echoBody map[string]interface{}

// Echo echoes whatever JSON a client sent through the API appending `"echoed": true`.
// If echoed is already set to true or the JSON payload is invalid it returns an error.
func Echo(b []byte) (echoBody, error) {
	if len(b) == 0 {
		return echoBody{}, errors.New("invalid json")
	}
	var m echoBody
	err := json.Unmarshal(b, &m)
	if err != nil {
		return echoBody{}, errors.New("invalid json")
	}
	// if the input contains "echoed" and it is true, return a bad request error
	if echoed, ok := m["echoed"]; ok {
		if echoed == "true" || echoed == true {
			return echoBody{}, errors.New("echoed already set")
		}
	}
	// otherwise set "echoed" to true
	m["echoed"] = true
	return m, nil
}
