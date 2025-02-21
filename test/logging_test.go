package simplelogger_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	slog "github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/simplelogger/examples"
)

var (
	sl           *slog.SimpleLoggerSettings
	listSettings []slog.OptionsManager

	err error
)

type InteractionDB struct{}

func NewInteractionDB() *InteractionDB {
	return &InteractionDB{}
}

func (idb *InteractionDB) Write(msgType, msg string) error {
	return errors.New("the database is not available")
}

func TestMain(m *testing.M) {
	opt := &examples.OptionForTest{}
	opt.SetNameMessageType("error")
	opt.SetMaxLogFileSize(1024)
	opt.SetPathDirectory("logs")
	opt.SetWritingStdout(true)
	opt.SetWritingFile(true)
	opt.SetWritingDB(true)
	listSettings = append(listSettings, opt)

	opt = &examples.OptionForTest{}
	opt.SetNameMessageType("info")
	opt.SetMaxLogFileSize(1024)
	opt.SetPathDirectory("logs")
	opt.SetWritingStdout(true)
	opt.SetWritingFile(true)
	opt.SetWritingDB(true)
	listSettings = append(listSettings, opt)

	opt = &examples.OptionForTest{}
	opt.SetNameMessageType("debug")
	opt.SetMaxLogFileSize(1024)
	opt.SetPathDirectory("logs")
	opt.SetWritingStdout(true)
	opt.SetWritingFile(false)
	opt.SetWritingDB(false)
	listSettings = append(listSettings, opt)

	opt = &examples.OptionForTest{}
	opt.SetNameMessageType("warning")
	opt.SetMaxLogFileSize(1024)
	opt.SetPathDirectory("logs")
	opt.SetWritingStdout(true)
	opt.SetWritingFile(true)
	opt.SetWritingDB(false)
	listSettings = append(listSettings, opt)

	opt = &examples.OptionForTest{}
	opt.SetNameMessageType("CRITICAL")
	opt.SetMaxLogFileSize(1024)
	opt.SetPathDirectory("logs")
	opt.SetWritingStdout(true)
	opt.SetWritingFile(true)
	opt.SetWritingDB(false)
	listSettings = append(listSettings, opt)

	opt = &examples.OptionForTest{}
	opt.SetNameMessageType("row_case")
	opt.SetMaxLogFileSize(1024)
	opt.SetPathDirectory("logs")
	opt.SetWritingStdout(true)
	opt.SetWritingFile(true)
	opt.SetWritingDB(false)
	listSettings = append(listSettings, opt)

	ctx, _ /*ctxClose*/ := context.WithCancel(context.Background())
	options := slog.CreateOptions(listSettings...)

	sl, err = slog.NewSimpleLogger(ctx, "simplelogger", options)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestLoggingImplementation(t *testing.T) {
	t.Run("Тест 1. Проверяем запись сообщения типа 'error' в лог-файл", func(t *testing.T) {
		ok := sl.Write("error", "my ERROR test message")
		assert.True(t, ok)
	})

	t.Run("Тест 2. Проверяем запись некоторого количества сообщений типа 'info' в лог-файл", func(t *testing.T) {
		ticker := time.NewTicker(1 * time.Second)

		var count int
		for {
			if count == 7 {
				ticker.Stop()

				break
			}

			count++

			<-ticker.C

			ok := sl.Write("info", fmt.Sprintf("%d. my INFORMATION test message", count))
			assert.True(t, ok)
		}
	})

	t.Run("Тест 3. Проверяем запись сообщения типа 'debug' в лог-файл", func(t *testing.T) {
		ok := sl.Write("debug", "my DEBUG test message")
		assert.False(t, ok)
	})

	t.Run("Тест 4. Проверяем запись сообщения типа 'warning' в лог-файл", func(t *testing.T) {
		ok := sl.Write("warning", "my WARNING test message")
		assert.True(t, ok)
	})

	t.Run("Тест 5. Проверяем запись сообщения типа 'critical' в лог-файл", func(t *testing.T) {
		ok := sl.Write("critical", "my CRITICAL test message")
		assert.True(t, ok)
	})

	t.Run("Тест 6. Проверяем запись сообщения типа 'row_case' в лог-файл", func(t *testing.T) {
		ok := sl.Write("row_case", "my ROW_CASE test message")
		assert.True(t, ok)
	})

	t.Run("Тест 7. Инициализация взаимодействия с БД", func(t *testing.T) {
		sl.SetDataBaseInteraction(NewInteractionDB())

		ok := sl.Write("error", "some description of the error")
		assert.True(t, ok)

		ok = sl.Write("info", "some description of the information message")
		assert.True(t, ok)
	})
}
