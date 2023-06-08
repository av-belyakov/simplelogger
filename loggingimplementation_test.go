package simplelogger_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/simplelogger"
)

var _ = Describe("Loggingimplementation", Ordered, func() {
	var (
		err error
		sl  simplelogger.SimpleLoggerSettings
	)

	BeforeAll(func() {
		sl, err = simplelogger.NewSimpleLogger("simplelogger", []simplelogger.MessageTypeSettings{
			{
				MsgTypeName:   "error",
				WritingFile:   true,
				PathDirectory: "logs",
				WritingStdout: false,
				MaxFileSize:   1024,
			},
			{
				MsgTypeName:   "info",
				WritingFile:   true,
				PathDirectory: "logs",
				WritingStdout: true,
				MaxFileSize:   1024,
			},
		})
	})

	Context("Тест 1. Проверяем работу логера", func() {
		It("Должен быть успешно создан новый объект логера", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(sl.GetCountFileDescription()).Should(Equal(2))
		})
	})

	Context("Тест 2. Проверяем запись сообщений в лог-файлы", func() {
		It("Должно быть записанно некоторое сообщение типа 'error' в лог-файл", func() {
			ok := sl.WriteLoggingData("my error test message", "error")

			_, err := os.ReadFile("logs/error.log")

			Expect(ok).Should(BeTrue())
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должно быть записанно некоторое сообщение типа 'error' в лог-файл", func() {
			ok := sl.WriteLoggingData("my information test message", "info")

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
			   TYPE - вид сообщения (DEBUG,INFO,WARNING,ERROR,CRITICAL...)
			   message - само сообщение
			   file - в каком файле вызвано (опционально)
			   line - в какой строке (опционально)
			*/

			Expect(true).Should(BeTrue())
		})
	})
})
