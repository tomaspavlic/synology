package synology

import (
	"net/url"
	"strconv"
)

// https://global.download.synology.com/download/Document/Software/DeveloperGuide/Os/DSM/All/enu/DSM_Login_Web_API_Guide_enu.pdf
// https://global.download.synology.com/download/Document/Software/DeveloperGuide/Package/FileStation/All/enu/Synology_File_Station_API_Guide.pdf

const (
	apiHost = "10.180.0.3:5001"
)

type Api struct {
	Name string
	Path string
}

func buildRequestUrl(path, name, method, host string, version int, params map[string]string) string {
	query := url.Values{}
	query.Set("api", name)
	query.Set("version", strconv.Itoa(version))
	query.Set("method", method)
	query.Set("query", "all")

	for k, v := range params {
		query.Set(k, v)
	}

	u := url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     "webapi/" + path,
		RawQuery: query.Encode(),
	}

	return u.String()
}
