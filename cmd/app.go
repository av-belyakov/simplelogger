package simplelogger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/av-belyakov/simplelogger/internal"
)

const DEFAULT_MAX_SIZE = 1000000

// NewSimpleLogger создает новый логер
func NewSimpleLogger(ctx context.Context, rootDir string, msgtsl []internal.MessageTypeSettings) (internal.SimpleLoggerSettings, error) {
	listType := [...]string{"INFO", "ERROR", "DEBUG", "WARNING", "CRITICAL"}
	logTypeIsExist := func(str string) bool {
		for _, v := range listType {
			if strings.ToUpper(str) == v {
				return true
			}
		}

		return false
	}

	getRootPath := func(rootDir string) (string, error) {
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

	sls := internal.SimpleLoggerSettings{RootDir: rootDir}
	mtd := map[string]internal.MessageTypeData{}
	if rootDir == "" {
		return sls, fmt.Errorf("the variable \"rootDir\" is not definitely")
	}

	rootPath, err := getRootPath(rootDir)
	if err != nil {
		return sls, err
	}

	sls.RootPath = rootPath

	go func() {
		<-ctx.Done()

		sls.ClosingFiles()
	}()

	for _, v := range msgtsl {
		if !v.WritingFile {
			continue
		}

		pd := v.PathDirectory
		if !strings.HasPrefix(pd, "/") {
			pd = path.Join(sls.RootPath, v.PathDirectory)
		}

		if !logTypeIsExist(v.MsgTypeName) {
			return sls, fmt.Errorf("the message type can only be one of the following types INFO, ERROR, DEBUG, WARNING, CRITICAL")
		}

		if _, err := os.ReadDir(pd); err != nil {
			if err := os.Mkdir(pd, 0777); err != nil {
				return sls, err
			}
		}

		msgTypeName := strings.ToLower(v.MsgTypeName)
		fullFileName := path.Join(pd, msgTypeName+".log")
		f, err := os.OpenFile(fullFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return sls, err
		}

		l := log.New(f, "", log.LstdFlags)
		if v.MsgTypeName == "error" {
			l.SetFlags(log.Lshortfile | log.LstdFlags)
		}

		if v.MaxFileSize < DEFAULT_MAX_SIZE {
			v.MaxFileSize = DEFAULT_MAX_SIZE
		}

		mtd[msgTypeName] = internal.MessageTypeData{
			MessageTypeSettings: internal.MessageTypeSettings{
				MsgTypeName:   v.MsgTypeName,
				WritingFile:   v.WritingFile,
				PathDirectory: pd,
				WritingStdout: v.WritingStdout,
				MaxFileSize:   v.MaxFileSize,
			},
			FileName:        fullFileName,
			FileDescription: f,
			LogDescription:  l,
		}
	}

	sls.ListMessageType = mtd

	return sls, nil
}
