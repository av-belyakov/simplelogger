package simplelogger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getRootPath(rootDir string) (string, error) {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	tmp := strings.Split(currentDir, "/")

	if tmp[len(tmp)-1] == rootDir {
		return currentDir, nil
	}

	var path string = ""
	for _, v := range tmp {
		path += v + "/"

		if v == rootDir {
			return path, nil
		}
	}

	return path, nil
}

func logTypeIsExist(str string) bool {
	listType := [...]string{"INFO", "ERROR", "DEBUG", "WARNING", "CRITICAL"}
	for _, v := range listType {
		if strings.ToUpper(str) == v {
			return true
		}
	}

	return false
}

func getColorTypeMsg(msgType string) string {
	switch msgType {
	case "INFO":
		return fmt.Sprintf("%vINF%v", AnsiBrightGreen, AnsiReset)

	case "ERROR":
		return fmt.Sprintf("%vERR%v", AnsiBrightRed, AnsiReset)

	case "DEBUG":
		return fmt.Sprintf("%vDBG%v", AnsiBrightWhite, AnsiReset)

	case "WARNING":
		return fmt.Sprintf("%vWARN%v", AnsiBrightYellow, AnsiReset)

	case "CRITICAL":
		return fmt.Sprintf("%vCRIT%v", AnsiBrightMagenta, AnsiReset)

	case "DBI":
		return fmt.Sprintf("%v%v ERR DataBaseInteractor %v", AnsiBackgroundRed, AnsiBrightBlack, AnsiReset)

	}

	return msgType
}
