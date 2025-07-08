package simplelogger

import (
	"log"
	"os"
)

// Options настройки типов сообщений
type Options struct {
	nameMessageType string //наименование типа сообщения
	pathDirectory   string //путь до директорий с лог-файлами, если в начале строки есть символ "/" то считается что директория создается от'корня' файловой системы, если символ "/" отсутствует в начале строки, то директория с логами будет создана в текущей директории приложения
	maxLogFileSize  int    //максимальный размер файла, в байтах, при достижении которого выполняется архивирование сообщения (не менее 1000000)
	writingToStdout bool   //писать ли сообщения данного типа в stdout
	writingToFile   bool   //писать ли сообщения данного типа в файл
	writingToDB     bool   //писать ли сообщения данного типа в БД
}

// MessageTypeData содержит данные по типу сообщений
type messageTypeData struct {
	logDescription  *log.Logger //дескриптор логов
	fileDescription *os.File    //файловый дескриптор для записи файлов
	fileName        string      //наименование лог-файла
	Options
}

// SimpleLoggerSettings содержит параметры SimpleLogger
type SimpleLoggerSettings struct {
	listMessageType     map[string]messageTypeData //список типов сообщений
	dataBaseInteraction DataBaseInteractor         //интерфейс для взаимодействия с БД
	rootDir             string                     //основная директория приложения
	rootPath            string                     //полный путь до директории приложения
}
