package helpers

import "encoding/json"

func MustMarshal(obj any) string {
	str, _ := json.Marshal(obj)
	return string(str)
}
