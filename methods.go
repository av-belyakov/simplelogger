package simplelogger

import (
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ansiReset         = "\033[0m"
	ansiBrightRed     = "\033[91m"
	ansiBrightGreen   = "\033[92m"
	ansiBrightWhite   = "\033[37m"
	ansiBrightYellow  = "\033[93m"
	ansiBrightMagenta = "\033[95m"
)

func (sls *SimpleLoggerSettings) GetCountFileDescription() int {
	var num int
	for _, v := range sls.ListMessageType {
		if v.WritingToFile {
			num++
		}
	}

	return num
}

func (sls *SimpleLoggerSettings) GetListTypeFiles() []string {
	list := make([]string, 0, len(sls.ListMessageType))
	for _, v := range sls.ListMessageType {
		list = append(list, v.MsgTypeName)
	}

	return list
}

func (sls *SimpleLoggerSettings) WriteLoggingData(str, typeLogFile string) bool {
	mt, ok := sls.ListMessageType[typeLogFile]
	if !ok {

		return false
	}

	//в консоль выводим только следующие типы сообщений: INFO, ERROR, DEBUG, WARNING, CRITICAL
	if mt.WritingToStdout && logTypeIsExist(typeLogFile) {
		//пишем в stdout
		tns := strings.Split(time.Now().String(), " ")
		dateTime := fmt.Sprintf("%s %s", tns[0], tns[1][:8])

		os.Stdout.Write([]byte(fmt.Sprintf("%s %s - %s - %s\n", dateTime, getColorTypeMsg(strings.ToUpper(typeLogFile)), sls.rootDir, str)))
	}

	if !mt.WritingToFile {
		return false
	}

	//пишем в файл
	mt.LogDescription.Println(str)

	fi, err := mt.FileDescription.Stat()
	if err != nil {
		return false
	}

	if fi.Size() <= int64(mt.MaxFileSize) {
		return true
	}

	mt.FileDescription.Close()
	//выполняем архивирование файла
	sls.compressFile(mt.FileName)

	_ = os.Remove(mt.FileName)
	f, err := os.OpenFile(mt.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return false
	}

	l := log.New(f, "", log.LstdFlags)
	if mt.MsgTypeName == "error" {
		l.SetFlags(log.Lshortfile | log.LstdFlags)
	}

	sls.ListMessageType[typeLogFile] = messageTypeData{
		Options:         mt.Options,
		FileName:        mt.FileName,
		FileDescription: f,
		LogDescription:  l,
	}

	return true
}

func (sls SimpleLoggerSettings) compressFile(tm string) {
	timeNowUnix := time.Now().Unix()
	fn := strings.Replace(tm, ".log", "_"+strconv.FormatInt(timeNowUnix, 10)+".gz", -1)

	fileIn, err := os.Create(fn)
	if err != nil {
		return
	}
	defer fileIn.Close()

	zw := gzip.NewWriter(fileIn)
	zw.Name = fn

	fileOut, err := os.ReadFile(tm)
	if err != nil {
		return
	}

	if _, err := zw.Write(fileOut); err != nil {
		return
	}

	_ = zw.Close()
}

func (sls *SimpleLoggerSettings) closingFiles() {
	for _, v := range sls.ListMessageType {
		v.FileDescription.Close()
	}
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
		return fmt.Sprintf("%vINF%v", ansiBrightGreen, ansiReset)

	case "ERROR":
		return fmt.Sprintf("%vERR%v", ansiBrightRed, ansiReset)

	case "DEBUG":
		return fmt.Sprintf("%vDBG%v", ansiBrightWhite, ansiReset)

	case "WARNING":
		return fmt.Sprintf("%vWARN%v", ansiBrightYellow, ansiReset)

	case "CRITICAL":
		return fmt.Sprintf("%vCRIT%v", ansiBrightMagenta, ansiReset)

	}

	return msgType
}
