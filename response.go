package synology

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorDetail struct {
	Code int
	Path string
}

type ResponseError struct {
	Code   int
	Errors []errorDetail
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("status code: %d", e.Code)
}

type response[T any] struct {
	Error   ResponseError
	Success bool
	Data    T
}

func readResponse[T any](resp *http.Response) (T, error) {
	var response response[T]
	var data T

	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return data, err
	}

	if !response.Success {
		return data, response.Error
	}
	data = response.Data

	return data, err
}
