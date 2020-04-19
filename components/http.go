package components

import (
	"net/url"
)

func MakeURL(addr, path string, params map[string]string) string {
	baseUrl, err := url.Parse(addr + path)
	if err != nil {
		logger.Error(
			"module", "http",
			"msg", "make URL error",
			"addr", addr,
			"path", path,
			"params", params,
		)
		return ""
	}

	if len(params) > 0 {
		baseUrl.Path += "?"
		urlParam := url.Values{}
		for k, v := range params {
			urlParam.Add(k, v)
		}

		// Add Query Parameters to the URL
		baseUrl.RawQuery = urlParam.Encode() // Escape Query Parameters
	}
	return baseUrl.String()
}
