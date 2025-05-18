package utils

import (
	"encoding/json"
	"strings"
)

func SearchInArray[T any](data []T, keySearch string, t T) (out []T, err error) {
	var matchedRecords []T
	for _, rec := range data {
		str, _ := json.Marshal(rec)
		strData := string(str)
		strData = strings.ToLower(strData)
		keySearch = strings.ToLower(keySearch)
		keySearch = strings.Trim(keySearch, " ")
		keySearch = strings.TrimSpace(keySearch)
		contains := strings.Contains(strData, keySearch)
		if !contains {
			continue
		}
		matchedRecords = append(matchedRecords, rec)
	}
	return matchedRecords, nil
}
