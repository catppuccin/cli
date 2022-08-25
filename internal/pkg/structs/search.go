package structs

import (
	"encoding/json"
)

func UnmarshalSearch(data []byte) (SearchRes, error) {
	var s SearchRes
	//s := SearchRes{}
	err := json.Unmarshal(data, &s)
	return s, err
}

type SearchRes struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
