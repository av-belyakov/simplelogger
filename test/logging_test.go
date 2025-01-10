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
)

var (
	sl  *slog.SimpleLoggerSettings
	err error

	listSettings []slog.Options = []slog.Options{
		{
			MsgTypeName:     "error",
			WritingToFile:   true,
			PathDirectory:   "logs",
			WritingToStdout: true,
			MaxFileSize:     1024,
			WritingToDB:     true,
		},
		{
			MsgTypeName:     "info",
			WritingToFile:   true,
			PathDirectory:   "logs",
			WritingToStdout: true,
			MaxFileSize:     1024,
			WritingToDB:     true,
		},
		{
			MsgTypeName:     "debug",
			WritingToFile:   false,
			PathDirectory:   "logs",
			WritingToStdout: true,
			MaxFileSize:     1024,
		},
		{
			MsgTypeName:     "warning",
			WritingToFile:   true,
			PathDirectory:   "logs",
			WritingToStdout: true,
			MaxFileSize:     1024,
		},
		{
			MsgTypeName:     "CRITICAL",
			WritingToFile:   true,
			PathDirectory:   "logs",
			WritingToStdout: true,
			MaxFileSize:     1024,
		},
		{
			MsgTypeName:     "row_case",
			WritingToFile:   true,
			PathDirectory:   "logs",
			WritingToStdout: true,
			MaxFileSize:     1024,
		},
	}
)

type InteractionDB struct{}

func NewInteractionDB() *InteractionDB {
	return &InteractionDB{}
}

func (idb *InteractionDB) Write(msgType, msg string) error {
	return errors.New("the database is not available")
}

func TestMain(m *testing.M) {
	ctx, _ /*ctxClose*/ := context.WithCancel(context.Background())
	sl, err = slog.NewSimpleLogger(ctx, "simplelogger", listSettings)
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
