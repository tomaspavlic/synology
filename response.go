package synology

import (
	"encoding/json"
	"fmt"
)

type errorDetail struct {
	Code int
	Path string
}

type responseError struct {
	Code   int
	Errors []errorDetail
}

func (e responseError) Error() string {
	return fmt.Sprintf("status code: %d", e.Code)
}

type response[T any] struct {
	Error   responseError
	Success bool
	Data    T
}

func unmarshal[T any](resp []byte) (T, error) {
	var response response[T]
	var data T

	err := json.Unmarshal(resp, &response)
	if err != nil {
		return data, err
	}

	if !response.Success {
		return data, response.Error
	}
	data = response.Data

	return data, err
}
