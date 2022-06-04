package synology

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	infoApiVersion = 1
	infoApiName    = "SYNO.API.Info"

	// default paths
	authPath  = "auth.cgi"
	queryPath = "query.cgi"
	entryPath = "entry.cgi"

	// default methods
	queryMethod = "query"
)

type ApiDetail struct {
	MaxVersion int
	MinVersion int
	Path       string
}

type ApiInfo map[string]ApiDetail

type parameters map[string]string

type SynologyCore struct {
	host       string
	httpClient *http.Client
	sid        string
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

// makeRequest builds request url from provided API information and makes http request to Synology API using JSON payload.
func (s *SynologyCore) makeRequest(path, name, method string, version int, params parameters) ([]byte, error) {

	if params == nil {
		params = make(parameters)
	}

	if s.sid != "" {
		params["_sid"] = s.sid
	}

	url := buildRequestUrl(path, name, method, s.host, version, params)
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := s.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// RetrieveApiInformation provides available API info.
func (s *SynologyCore) ApiInfo() (ApiInfo, error) {
	response, err := s.makeRequest(
		queryPath,
		infoApiName,
		queryMethod,
		infoApiVersion,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return unmarshal[ApiInfo](response)
}
