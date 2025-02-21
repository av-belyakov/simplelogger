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

// SetDataBaseInteraction устанавливает метод для взаимодействия с БД
func (sls *SimpleLoggerSettings) SetDataBaseInteraction(dbi DataBaseInteractor) {
	sls.dataBaseInteraction = dbi
}

// GetCountFileDescription количество открытых файловых дескрипторов
func (sls *SimpleLoggerSettings) GetCountFileDescription() int {
	var num int
	for _, v := range sls.listMessageType {
		if v.writingToFile {
			num++
		}
	}

	return num
}

// GetListTypeFiles список типов файлов
func (sls *SimpleLoggerSettings) GetListTypeFiles() []string {
	list := make([]string, 0, len(sls.listMessageType))
	for _, v := range sls.listMessageType {
		list = append(list, v.nameMessageType)
	}

	return list
}

// WriteData запись логов (в сигнатуре функции сначала идет ТИП сообщения, затем САМО сообщение)
func (sls *SimpleLoggerSettings) Write(typeLog, msg string) bool {
	if sls.listMessageType == nil {
		return false
	}

	mt, ok := sls.listMessageType[typeLog]
	if !ok {
		return false
	}

	tns := strings.Split(time.Now().String(), " ")
	dateTime := fmt.Sprintf("%s %s", tns[0], tns[1][:8])

	//в консоль выводим только следующие типы сообщений: INFO, ERROR, DEBUG, WARNING, CRITICAL
	if mt.writingToStdout && logTypeIsExist(typeLog) {
		os.Stdout.Write([]byte(fmt.Sprintf("%s %s - %s - %s\n", dateTime, getColorTypeMsg(strings.ToUpper(typeLog)), sls.rootDir, msg)))
	}

	//запись сообщений в БД
	if mt.writingToDB && sls.dataBaseInteraction != nil {
		if err := sls.dataBaseInteraction.Write(typeLog, msg); err != nil {
			os.Stdout.Write([]byte(fmt.Sprintf("%s %s - %s - %s\n", dateTime, getColorTypeMsg("DBI"), sls.rootDir, err.Error())))
		}
	}

	if !mt.writingToFile {
		return false
	}

	//пишем в файл
	mt.logDescription.Println(msg)

	fi, err := mt.fileDescription.Stat()
	if err != nil {
		return false
	}

	if fi.Size() <= int64(mt.maxLogFileSize) {
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
	if mt.nameMessageType == "error" {
		l.SetFlags(log.Lshortfile | log.LstdFlags)
	}

	sls.listMessageType[typeLog] = messageTypeData{
		Options:         mt.Options,
		fileName:        mt.fileName,
		fileDescription: f,
		logDescription:  l,
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
	for _, v := range sls.listMessageType {
		v.fileDescription.Close()
	}
}
