package util

import (
	"github.com/cihub/seelog"
	"os"
	"fmt"
	"path/filepath"
)

const CONF_FILE_NAME = "logger.xml"

func InitLogger() {
	path := GetConfPathFromOsArgs()
	file := filepath.Join(path, CONF_FILE_NAME)
	fmt.Printf("conf file: %s\n", file)
	logger, err := seelog.LoggerFromConfigAsFile(file)
	if err != nil {
		seelog.Critical("err parsing config log file ", err)
		os.Exit(1)
	}

	seelog.ReplaceLogger(logger)
}
