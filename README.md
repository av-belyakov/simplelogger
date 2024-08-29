#README

Модуль "Simpleloger" осуществляет запись ошибок, информационных сообщений и любых других сообщений в лог-файлы.

Типы записываемых сообщений:

- INFO
- ERROR
- DEBUG
- WARNING
- CRITICAL

Лог-файлы могут писатся как в корневую директорию проекта, для этого надо указать папку
для записи логов, например, папка 'logs', так и в любую другую директорию, нужно написать полный путь
до этой директории. Создание файла будет выполненно в указанной директории если для этого будет достаточно прав.

Модуль выполняет архивирование старых лог-файлов в формат gzip.

Модуль может отправлять ошибки или другие сообщения на вывод stdout.

Модуль предназначен только для ОС Unix.
