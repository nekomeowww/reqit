package reqit

import "net/url"

func mapToQuery(param map[string]string) url.Values {
	queries := url.Values{}
	for key, value := range param {
		queries.Add(key, value)
	}

	return queries
}
