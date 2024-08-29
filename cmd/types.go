package simplelogger

import (
	"log"
	"os"
)

// messageTypeSettings настройки типов сообщений
// WritingFile - писать ли сообщения данного типа в файл
// WritingStdout - писать ли сообщения данного типа в stdout
// MsgTypeName - наименование типа сообщения
// MaxFileSize - максимальный размер файла, в байтах, при достижении которого выполняется архивирование сообщения (не менее 1000000)
// PathDirectory - путь до директорий с лог-файлами, если в начале строки есть символ "/" то считается что директория создается от
// 'корня' файловой системы, если символ "/" отсутствует в начале строки, то директория с логами будет создана в текущей директории
// приложения
type messageTypeSettings struct {
	WritingFile   bool
	WritingStdout bool
	MaxFileSize   int
	MsgTypeName   string
	PathDirectory string
}

// MessageTypeData содержит данные по типу сообщений
// FileName - наименование лог-файла
// FileDescription - файловый дескриптор для записи файлов
// LogDescription - дескриптор логов
type messageTypeData struct {
	messageTypeSettings
	FileName        string
	FileDescription *os.File
	LogDescription  *log.Logger
}

// SimpleLoggerSettings содержит параметры SimpleLogger
// RootDir - основная директория приложения
// RootPath - полный путь до директории приложения
// ListMessageType - список типов сообщений
type SimpleLoggerSettings struct {
	rootDir         string
	rootPath        string
	ListMessageType map[string]messageTypeData
}
