package internal

import (
	"log"
	"os"
)

// MessageTypeSettings настройки типов сообщений
// WritingFile - писать ли сообщения данного типа в файл
// WritingStdout - писать ли сообщения данного типа в stdout
// MsgTypeName - наименование типа сообщения
// MaxFileSize - максимальный размер файла, в байтах, при достижении которого выполняется архивирование сообщения (не менее 1000000)
// PathDirectory - путь до директорий с лог-файлами, если в начале строки есть символ "/" то считается что директория создается от
// 'корня' файловой системы, если символ "/" отсутствует в начале строки, то директория с логами будет создана в текущей директории
// приложения
type MessageTypeSettings struct {
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
type MessageTypeData struct {
	MessageTypeSettings
	FileName        string
	FileDescription *os.File
	LogDescription  *log.Logger
}

// SimpleLoggerSettings содержит параметры SimpleLogger
// RootDir - основная директория приложения
// RootPath - полный путь до директории приложения
// ListMessageType - список типов сообщений
type SimpleLoggerSettings struct {
	RootDir         string
	RootPath        string
	ListMessageType map[string]MessageTypeData
}
