package convert

import (
	"gitote.com/gitote/gitote/pkg/setting"
)

// ToCorrectPageSize makes sure page size is in allowed range.
func ToCorrectPageSize(size int) int {
	if size <= 0 {
		size = 10
	} else if size > setting.API.MaxResponseItems {
		size = setting.API.MaxResponseItems
	}
	return size
}
