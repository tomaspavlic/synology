package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

const (
	apiInfoApiName = "SYNO.API.Info"
	queryPath      = "query.cgi"
	queryMethod    = "query"
)

type ApiDetail struct {
	MaxVersion int
	MinVersion int
	Path       string
}

type ApiInfo map[string]ApiDetail

type SynologyCore struct {
	info       ApiInfo
	host       string
	httpClient *http.Client
	authCookie *http.Cookie
}

func NewSynologyCore(host string, port int) *SynologyCore {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &SynologyCore{
		httpClient: client,
		host:       fmt.Sprintf("%s:%d", host, port),
	}
}

func (s *SynologyCore) makeRequest(path, name, method string, version int, params map[string]string) (*http.Response, error) {
	url := buildRequestUrl(path, name, method, s.host, version, params)
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	// include authentication cookie if exists
	if s.authCookie != nil {
		request.AddCookie(s.authCookie)
	}
	resp, err := s.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SynologyCore) RetrieveApiInformation() (ApiInfo, error) {
	response, err := s.makeRequest(queryPath, apiInfoApiName, queryMethod, 1, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := readResponse[ApiInfo](response)

	return data, err
}

func (s *SynologyCore) Find(apiName string) (*Api, error) {
	// make sure all information is first loaded from API
	if s.info == nil {
		info, err := s.RetrieveApiInformation()

		if err != nil {
			return nil, err
		}

		s.info = info
	}

	if api, ok := s.info[apiName]; ok {
		return &Api{
			Name: apiName,
			Path: api.Path,
		}, nil
	}

	return nil, fmt.Errorf("provided api name '%s' was not found", apiName)
}
