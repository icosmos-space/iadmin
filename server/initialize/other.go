package initialize

import (
	"bufio"
	"os"
	"strings"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/utils"
)

func OtherInit() {
	if strings.TrimSpace(global.IADMIN_CONFIG.JWT.ExpiresTime) == "" {
		global.IADMIN_CONFIG.JWT.ExpiresTime = "7d"
	}
	if strings.TrimSpace(global.IADMIN_CONFIG.JWT.BufferTime) == "" {
		global.IADMIN_CONFIG.JWT.BufferTime = "1d"
	}

	dr, err := utils.ParseDuration(global.IADMIN_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.IADMIN_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
	file, err := os.Open("go.mod")
	if err == nil && global.IADMIN_CONFIG.AutoCode.Module == "" {
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		global.IADMIN_CONFIG.AutoCode.Module = strings.TrimPrefix(scanner.Text(), "module ")
	}
}
