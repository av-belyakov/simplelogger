package simplelogger

// DataBaseInteractor реализует методы для взаимодействия модуля с БД
type DataBaseInteractor interface {
	Write(msgType, msg string) error
}
