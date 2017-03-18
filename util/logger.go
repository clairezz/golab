package util

import (
	"github.com/cihub/seelog"
	"os"
)

const CONFIG_FILE  = "D:\\www\\godev\\src\\golab\\conf\\logger.xml"

// TODO 配置文件相对路径
func InitLogger() {
	logger, err := seelog.LoggerFromConfigAsFile(CONFIG_FILE)
	if err != nil {
		seelog.Critical("err parsing config log file ", err)
		os.Exit(1)
	}

	seelog.ReplaceLogger(logger)
}
