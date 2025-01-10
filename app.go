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
		msgTypeName := strings.ToLower(v.MsgTypeName)

		pd := v.PathDirectory
		if !strings.HasPrefix(pd, "/") {
			pd = path.Join(sls.rootPath, v.PathDirectory)
		}

		maxFileSize := DEFAULT_MAX_SIZE
		if v.MaxFileSize > 1000 {
			maxFileSize = v.MaxFileSize
		}

		mtd[msgTypeName] = messageTypeData{
			Options: Options{
				WritingToFile:   v.WritingToFile,
				WritingToStdout: v.WritingToStdout,
				MaxFileSize:     maxFileSize,
				MsgTypeName:     v.MsgTypeName,
				PathDirectory:   pd,
				WritingToDB:     v.WritingToDB,
			}}

		if !v.WritingToFile {
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
		if v.MsgTypeName == "error" {
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
