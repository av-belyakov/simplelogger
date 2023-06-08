package simplelogger_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"simplelogger"
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
})
