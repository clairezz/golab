package util

import (
	"os"
	"path/filepath"
	"flag"
	"os/exec"
	"fmt"
)

func GetConfPathFromOsArgs() string {
	file, _ := exec.LookPath(os.Args[0])
	fmt.Println("file: ", file)

	path, _ := filepath.Abs(file)
	fmt.Println("path: ", path)

	root := filepath.Dir(path)
	fmt.Println("root: ", root)

	var configPath string
	flag.StringVar(&configPath, "conf", "", "")
	flag.Parse()

	if "" == configPath {
		configPath = filepath.Join(root, "conf")
	} else {
		configPath, _ = filepath.Abs(configPath)
	}

	return configPath
}
