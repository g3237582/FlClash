package anet

import (
	"os"
	"runtime"
	"strconv"
)

func IsForceAnet() bool {
	return forceAnet(runtime.GOOS, os.Getenv("FORCE_ANET"))
}

func forceAnet(goos, env string) bool {
	if goos == "android" {
		return true
	}
	force, _ := strconv.ParseBool(env)
	return force
}
