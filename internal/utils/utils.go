package utils

import "strings"

func IsValueExist(data map[string]string, searchText string) (string, bool) {
	if data != nil {
		for key, value := range data {
			if strings.EqualFold(searchText, value) {
				return key, true
			}
		}
	}
	return "", false
}
