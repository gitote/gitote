// +build miniwinsvc

package setting

import (
	_ "gitote.com/gitote/minwinsvc"
)

func init() {
	SupportMiniWinService = true
}
