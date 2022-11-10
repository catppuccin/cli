package structs

import (
	"encoding/json"
)

func UnmarshalSearch(data []byte) (SearchRes, error) {
	var s SearchRes
	// s := SearchRes{}
	err := json.Unmarshal(data, &s)
	return s, err
}

type SearchRes []SearchEntry

type SearchEntry struct {
	Name   string   `json:"name"`
	Stars  int      `json:"stars"`
	Topics []string `json:"topics"`
}
