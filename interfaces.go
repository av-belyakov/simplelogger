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

/*
type LogSet struct {
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	WritingDB     bool   `validate:"required" yaml:"writingDB"`
}


			WritingToStdout: v.WritingStdout,
			WritingToFile:   v.WritingFile,
			WritingToDB:     v.WritingDB,
			MsgTypeName:     v.MsgTypeName,
			PathDirectory:   v.PathDirectory,
			MaxFileSize:     v.MaxFileSize,
*/
