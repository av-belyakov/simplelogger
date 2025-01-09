package simplelogger

import (
	"log"
	"os"
)

// Options настройки типов сообщений
type Options struct {
	WritingToDB     bool   //писать ли сообщения данного типа в БД
	WritingToFile   bool   //писать ли сообщения данного типа в файл
	WritingToStdout bool   //писать ли сообщения данного типа в stdout
	MaxFileSize     int    //максимальный размер файла, в байтах, при достижении которого выполняется архивирование сообщения (не менее 1000000)
	MsgTypeName     string //наименование типа сообщения
	PathDirectory   string //путь до директорий с лог-файлами, если в начале строки есть символ "/" то считается что директория создается от'корня' файловой системы, если символ "/" отсутствует в начале строки, то директория с логами будет создана в текущей директории приложения
}

// MessageTypeData содержит данные по типу сообщений
type messageTypeData struct {
	Options
	FileName        string      //наименование лог-файла
	FileDescription *os.File    //файловый дескриптор для записи файлов
	LogDescription  *log.Logger //дескриптор логов
}

// SimpleLoggerSettings содержит параметры SimpleLogger
type SimpleLoggerSettings struct {
	rootDir             string                     //основная директория приложения
	rootPath            string                     //полный путь до директории приложения
	dataBaseInteraction DataBaseInteractor         //интерфейс для взаимодействия с БД
	ListMessageType     map[string]messageTypeData //список типов сообщений
}
