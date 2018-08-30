package gitote

import (
	"net/http"
)

var jsonHeader = http.Header{"content-type": []string{"application/json"}}

func Bool(v bool) *bool {
	return &v
}

func String(v string) *string {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}
