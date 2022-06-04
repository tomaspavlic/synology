package synology

import (
	"net/url"
	"strconv"
)

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
