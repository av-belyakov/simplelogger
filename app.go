package simplelogger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Options настройки типов сообщений
// WritingToFile - писать ли сообщения данного типа в файл
// WritingToStdout - писать ли сообщения данного типа в stdout
// MsgTypeName - наименование типа сообщения
// MaxFileSize - максимальный размер файла, в байтах, при достижении которого выполняется архивирование сообщения (не менее 1000000)
// PathDirectory - путь до директорий с лог-файлами, если в начале строки есть символ "/" то считается что директория создается от
// 'корня' файловой системы, если символ "/" отсутствует в начале строки, то директория с логами будет создана в текущей директории
// приложения
type Options struct {
	WritingToFile   bool
	WritingToStdout bool
	MaxFileSize     int
	MsgTypeName     string
	PathDirectory   string
}

// NewSimpleLogger создает новый логер
func NewSimpleLogger(ctx context.Context, rootDir string, opt []Options) (*SimpleLoggerSettings, error) {
	const DEFAULT_MAX_SIZE = 1000000

	/*
		listType := [...]string{"INFO", "ERROR", "DEBUG", "WARNING", "CRITICAL"}
		logTypeIsExist := func(str string) bool {
			for _, v := range listType {
				if strings.ToUpper(str) == v {
					return true
				}
			}

			return false
		}
	*/

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

	sls := SimpleLoggerSettings{rootDir: rootDir}
	mtd := map[string]messageTypeData{}
	if rootDir == "" {
		return &sls, fmt.Errorf("the variable \"rootDir\" is not definitely")
	}

	rootPath, err := getRootPath(rootDir)
	if err != nil {
		return &sls, err
	}

	sls.rootPath = rootPath

	go func() {
		<-ctx.Done()

		sls.closingFiles()
	}()

	for _, v := range opt {
		if !v.WritingToFile {
			continue
		}

		pd := v.PathDirectory
		if !strings.HasPrefix(pd, "/") {
			pd = path.Join(sls.rootPath, v.PathDirectory)
		}

		//if !logTypeIsExist(v.MsgTypeName) {
		//	return &sls, fmt.Errorf("the message type can only be one of the following types INFO, ERROR, DEBUG, WARNING, CRITICAL")
		//}

		if _, err := os.ReadDir(pd); err != nil {
			if err := os.Mkdir(pd, 0777); err != nil {
				return &sls, err
			}
		}

		msgTypeName := strings.ToLower(v.MsgTypeName)
		fullFileName := path.Join(pd, msgTypeName+".log")
		f, err := os.OpenFile(fullFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return &sls, err
		}

		l := log.New(f, "", log.LstdFlags)
		if v.MsgTypeName == "error" {
			l.SetFlags(log.Lshortfile | log.LstdFlags)
		}

		if v.MaxFileSize < DEFAULT_MAX_SIZE {
			v.MaxFileSize = DEFAULT_MAX_SIZE
		}

		mtd[msgTypeName] = messageTypeData{
			messageTypeSettings: messageTypeSettings{
				WritingFile:   v.WritingToFile,
				WritingStdout: v.WritingToStdout,
				MaxFileSize:   v.MaxFileSize,
				MsgTypeName:   v.MsgTypeName,
				PathDirectory: pd,
			},
			FileName:        fullFileName,
			FileDescription: f,
			LogDescription:  l,
		}
	}

	sls.ListMessageType = mtd

	return &sls, nil
}
