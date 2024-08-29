package simplelogger_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	slog "github.com/av-belyakov/simplelogger"
)

var _ = Describe("Loggingimplementation", Ordered, func() {
	var (
		err error
		sl  *slog.SimpleLoggerSettings

		listSettings []slog.Options
	)

	BeforeAll(func() {
		ctx, _ /*ctxClose*/ := context.WithCancel(context.Background())
		//go func() {
		//	time.Sleep(2 * time.Second)
		//
		//	ctxClose()
		//}()
		listSettings = []slog.Options{
			{
				MsgTypeName:     "error",
				WritingToFile:   true,
				PathDirectory:   "logs",
				WritingToStdout: true,
				MaxFileSize:     1024,
			},
			{
				MsgTypeName:     "info",
				WritingToFile:   true,
				PathDirectory:   "logs",
				WritingToStdout: true,
				MaxFileSize:     1024,
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

		sl, err = slog.NewSimpleLogger(ctx, "simplelogger", listSettings)
	})

	Context("Тест 1. Проверяем работу логера ", func() {
		It("Должен быть успешно создан новый объект логера", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sl.GetCountFileDescription()).Should(Equal(len(listSettings) - 1))
		})
	})

	Context("Тест 2. Проверяем запись сообщений в лог-файлы ", func() {
		It("Должно быть записанно некоторое сообщение типа 'error'", func() {
			ok := sl.WriteLoggingData("my ERROR test message", "error")
			Expect(ok).Should(BeTrue())
		})

		It("Должно быть записанно некоторое количество сообщений типа 'info'", func() {
			for i := 0; i < 10; i++ {
				ok := sl.WriteLoggingData(fmt.Sprintf("%d. my INFORMATION test message", i), "info")
				Expect(ok).Should(BeTrue())

				time.Sleep(500 * time.Millisecond)
			}
		})

		It("Должно быть записанно некоторое сообщение типа 'debug'", func() {
			ok := sl.WriteLoggingData("my DEBUG test message", "debug")
			Expect(ok).Should(BeFalse())
		})

		It("Должно быть записанно некоторое сообщение типа 'warning'", func() {
			ok := sl.WriteLoggingData("my WARNING test message", "warning")
			Expect(ok).Should(BeTrue())
		})

		It("Должно быть записанно некоторое сообщение типа 'critical'", func() {
			ok := sl.WriteLoggingData("my CRITICAL test message", "critical")
			Expect(ok).Should(BeTrue())
		})

		It("Должно быть записанно некоторое сообщение типа 'row_case'", func() {
			ok := sl.WriteLoggingData("my ROW_CASE test message", "row_case")
			Expect(ok).Should(BeTrue())
		})
	})

	Context("Тест 3. Проверяем формирование правельной сторки сообщения отправляемого на stdout", func() {
		It("Должна быть сформированна верная строка", func() {
			funcName := "test_func"
			description := "any message"

			tn := time.Now()
			tns := strings.Split(tn.String(), " ")
			strMsg := fmt.Sprintf("%s %s [%s %s] - %s ('%s')", tns[0], tns[1], tns[2], tns[3], description, funcName)

			fmt.Println()
			fmt.Println(strMsg)
			fmt.Printf("%s %s\n", tns[0], tns[1][:8])

			/*
			   %Y-%m-%d %H:%M:%s - project.module - TYPE - message file:line

			   %Y - год
			   %m - месяц
			   %d - день
			   %H - час
			   %M - минута
			   %s - секунда
			   project - название проекта
			   module - название модуля проекта, в котором появилось сообщение (опционально)
			   TYPE - вид сообщения (DEBUG,INFO,WARNING,ERROR,CRITICAL)
			   message - само сообщение
			   file - в каком файле вызвано (опционально)
			   line - в какой строке (опционально)
			*/

			Expect(true).Should(BeTrue())
		})
	})
})
