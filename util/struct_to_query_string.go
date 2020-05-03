package util

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// StructToQueryString converts flat struct to query string
func StructToQueryString(in interface{}) (out string, err error) {
	jsonBytes, err := json.Marshal(in)
	if err != nil {
		return
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return
	}

	vals := url.Values{}
	for key, v := range data {
		vSlice, ok := v.([]interface{})
		if ok {
			for _, s := range vSlice {
				vals.Add(key, fmt.Sprintf("%v", s))
			}
			continue
		}
		vals.Add(key, fmt.Sprintf("%v", v))
	}

	return vals.Encode(), nil
}
