package synology

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
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
func (s *SynologyCore) makeRequest(path, name, method string, version int, params map[string]string) ([]byte, error) {

	if params == nil {
		params = make(map[string]string)
	}

	// include authentication cookie if exists
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
func (s *SynologyCore) RetrieveApiInformation() (ApiInfo, error) {
	response, err := s.makeRequest(queryPath, apiInfoApiName, queryMethod, 1, nil)
	if err != nil {
		return nil, err
	}

	return unmarshal[ApiInfo](response)
}

// Find finds provided API name from available API info.
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

func (s *SynologyCore) FileStation() (*FileStation, error) {
	listApi, err := s.Find(fileStationApiName)

	fs := &FileStation{
		core:    s,
		listApi: listApi,
	}

	return fs, err
}

// SurveillanceStation tries to get API for SurveillanceStation.
func (s *SynologyCore) SurveillanceStation() (*SurveillanceStation, error) {
	cameraApi, err := s.Find(surveillanceStationCameraApiName)

	ss := &SurveillanceStation{
		core:      s,
		cameraApi: cameraApi,
	}

	return ss, err
}
