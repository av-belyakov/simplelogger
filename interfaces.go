package simplelogger

// DataBaseInteractor реализует методы для взаимодействия модуля с БД
type DataBaseInteractor interface {
	Write(msgType, msg string) error
}

type OptionsManager interface {
	OptionsSetter
	OptionsGetter
}

type OptionsSetter interface {
	SetNameMessageType(string) error
	SetMaxLogFileSize(int) error
	SetPathDirectory(string) error
	SetWritingStdout(bool)
	SetWritingFile(bool)
	SetWritingDB(bool)
}

type OptionsGetter interface {
	GetNameMessageType() string
	GetMaxLogFileSize() int
	GetPathDirectory() string
	GetWritingStdout() bool
	GetWritingFile() bool
	GetWritingDB() bool
}
