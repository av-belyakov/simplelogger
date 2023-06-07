package simplelogger

import (
	"compress/gzip"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// MessageTypeSettings настройки типов сообщений
// MsgTypeName - наименование типа сообщения
// WritingFile - писать ли сообщения данного типа в файл
// PathDirectory - путь до директорий с лог-файлами (если путь пустой то используется наименование директории приложения)
// DirectoryName - наименование директории в которой сохраняется лог-файлы с сообщениями этого типа
// WritingStdout - писать ли сообщения данного типа в stdout
// MaxFileSize - максимальный размер файла при достижении которого выполняется архивирование сообщения (не менее 1000000)
type MessageTypeSettings struct {
	MsgTypeName   string
	WritingFile   bool
	PathDirectory string
	DirectoryName string
	WritingStdout bool
	MaxFileSize   int
}

// MessageTypeData содержит данные по типу сообщений
// fileName - наименование лог-файла
// fileDescription - файловый дескриптор для записи файлов
// logDescription - дескриптор логов
type MessageTypeData struct {
	MessageTypeSettings
	fileName        string
	fileDescription *os.File
	logDescription  *log.Logger
}

// SimpleLoggerSettings содержит параметры SimpleLogger
// ootPath - основная директория приложения
// listMessageType - список типов сообщений
type SimpleLoggerSettings struct {
	rootPath        string
	listMessageType map[string]MessageTypeData
}

// NewSimpleLogger создает новый логер
func NewSimpleLogger(rootDir string, msgtsl []MessageTypeSettings) (SimpleLoggerSettings, error) {
	const DEFAULT_MAX_SIZE = 1000000

	// ТОЛЬКО ДЛЯ ТЕСТА
	//const DEFAULT_MAX_SIZE = 1000

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

	sls := SimpleLoggerSettings{}
	mtd := map[string]MessageTypeData{}
	if rootDir == "" {
		return sls, fmt.Errorf("the variable 'rootDir' is not definitely")
	}

	rootPath, err := getRootPath(rootDir)
	if err != nil {
		return sls, err
	}

	sls.rootPath = rootPath

	for _, v := range msgtsl {
		pd := path.Join(v.PathDirectory, "/", v.DirectoryName)
		if v.PathDirectory == "" {
			pd = path.Join(sls.rootPath, v.DirectoryName)
		}

		if v.WritingFile {
			if _, err := os.ReadDir(pd); err != nil {
				if err := os.Mkdir(pd, 0777); err != nil {
					return sls, errors.New("error: it is not possible to create a directory for log files")
				}
			}

			fullfn := path.Join(pd, v.MsgTypeName+".log")
			f, err := os.OpenFile(fullfn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return sls, fmt.Errorf("error: it is impossible to create a log file %s", fullfn)
			}

			l := log.New(f, "", log.LstdFlags)
			if v.MsgTypeName == "error" {
				l.SetFlags(log.Lshortfile | log.LstdFlags)
			}

			if v.MaxFileSize < DEFAULT_MAX_SIZE {
				v.MaxFileSize = DEFAULT_MAX_SIZE
			}

			mtd[v.MsgTypeName] = MessageTypeData{
				MessageTypeSettings: MessageTypeSettings{
					MsgTypeName:   v.MsgTypeName,
					WritingFile:   v.WritingFile,
					PathDirectory: v.PathDirectory,
					DirectoryName: v.PathDirectory,
					WritingStdout: v.WritingStdout,
					MaxFileSize:   v.MaxFileSize,
				},
				fileName:        fullfn,
				fileDescription: f,
				logDescription:  l,
			}
		}
	}

	sls.listMessageType = mtd

	return sls, nil
}

func (sls *SimpleLoggerSettings) ClosingFiles() {
	for _, v := range sls.listMessageType {
		v.fileDescription.Close()
	}
}

func (sls *SimpleLoggerSettings) GetCountFileDescription() int {
	var num int
	for _, v := range sls.listMessageType {
		if v.WritingFile {
			num++
		}
	}

	return num
}

func (sls *SimpleLoggerSettings) GetListTypeFiles() []string {
	list := make([]string, 0, len(sls.listMessageType))
	for _, v := range sls.listMessageType {
		list = append(list, v.MsgTypeName)
	}

	return list
}

func (sls *SimpleLoggerSettings) WriteLoggingData(str, typeLogFile string) bool {
	mt, ok := sls.listMessageType[typeLogFile]
	if !ok {
		return false
	}

	if mt.WritingStdout {
		//пишем в stdout
		os.Stdout.Write([]byte(str))
	}

	if mt.WritingFile {
		//пишем в файл
		mt.logDescription.Println(str)
	}

	fi, err := mt.fileDescription.Stat()
	if err != nil {
		return false
	}

	if fi.Size() <= int64(mt.MaxFileSize) {
		return true
	}

	mt.fileDescription.Close()
	//выполняем архивирование файла
	sls.compressFile(mt.fileName)

	_ = os.Remove(mt.fileName)
	f, err := os.OpenFile(mt.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return false
	}

	l := log.New(f, "", log.LstdFlags)
	if mt.MsgTypeName == "error" {
		l.SetFlags(log.Lshortfile | log.LstdFlags)
	}

	sls.listMessageType[typeLogFile] = MessageTypeData{
		MessageTypeSettings: mt.MessageTypeSettings,
		fileName:            mt.fileName,
		fileDescription:     f,
		logDescription:      l,
	}

	return true
}

func (sls SimpleLoggerSettings) compressFile(tm string) {
	timeNowUnix := time.Now().Unix()
	fn := strings.Replace(tm, ".log", "_"+strconv.FormatInt(timeNowUnix, 10)+".gz", -1)

	fmt.Println("func 'compressFile', tm = ", tm, " fn = ", fn)

	//fileIn, err := os.Create(path.Join(cilf.locationLogDirectory, cilf.nameLogDirectory, fn))
	fileIn, err := os.Create(fn)
	if err != nil {

		fmt.Println("func 'compressFile', 111 ERROR = ", err)

		return
	}
	defer fileIn.Close()

	zw := gzip.NewWriter(fileIn)
	zw.Name = fn

	//fileOut, err := ioutil.ReadFile(path.Join(cilf.locationLogDirectory, cilf.nameLogDirectory, cilf.fileNameType[tm]))
	//fileOut, err := ioutil.ReadFile(tm)
	fileOut, err := os.ReadFile(tm)
	if err != nil {

		fmt.Println("func 'compressFile', 222 ERROR = ", err)

		return
	}

	if _, err := zw.Write(fileOut); err != nil {

		fmt.Println("func 'compressFile', 333 ERROR = ", err)

		return
	}

	_ = zw.Close()
}
