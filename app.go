package simplelogger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

// NewSimpleLogger создает новый логер
func NewSimpleLogger(ctx context.Context, rootDir string, opt []Options) (*SimpleLoggerSettings, error) {
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
		msgTypeName := strings.ToLower(v.nameMessageType)

		pd := v.pathDirectory
		if !strings.HasPrefix(pd, "/") {
			pd = path.Join(sls.rootPath, v.pathDirectory)
		}

		maxFileSize := DEFAULT_MAX_SIZE
		if v.maxLogFileSize > 1000 {
			maxFileSize = v.maxLogFileSize
		}

		mtd[msgTypeName] = messageTypeData{
			Options: Options{
				writingToDB:     v.writingToDB,
				writingToFile:   v.writingToFile,
				writingToStdout: v.writingToStdout,
				maxLogFileSize:  maxFileSize,
				nameMessageType: v.nameMessageType,
				pathDirectory:   pd,
			}}

		if !v.writingToFile {
			continue
		}

		if _, err := os.ReadDir(pd); err != nil {
			if err := os.Mkdir(pd, 0777); err != nil {
				return &sls, err
			}
		}

		fullFileName := path.Join(pd, msgTypeName+".log")
		f, err := os.OpenFile(fullFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return &sls, err
		}

		l := log.New(f, "", log.LstdFlags)
		if v.nameMessageType == "error" {
			l.SetFlags(log.Lshortfile | log.LstdFlags)
		}

		if data, ok := mtd[msgTypeName]; ok {
			data.FileName = fullFileName
			data.FileDescription = f
			data.LogDescription = l

			mtd[msgTypeName] = data
		}
	}

	sls.ListMessageType = mtd

	return &sls, nil
}

// CreateOptions создает список опций логирования
func CreateOptions(rawOpt ...OptionsManager) []Options {
	options := make([]Options, 0, len(rawOpt))

	for _, v := range rawOpt {
		options = append(options, Options{
			writingToStdout: v.GetWritingStdout(),
			writingToFile:   v.GetWritingFile(),
			writingToDB:     v.GetWritingDB(),
			nameMessageType: v.GetNameMessageType(),
			pathDirectory:   v.GetPathDirectory(),
			maxLogFileSize:  v.GetMaxLogFileSize(),
		})
	}

	return options
}
