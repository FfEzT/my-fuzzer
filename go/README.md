### запуск программы
`go run cmd/fuzzer/main.go -u HOSTNAME -w /path/to/wordlist -d "user=admin&pass=FUZZ"`

`go run cmd/fuzzer/main.go --help`

-c - это кол-во воркеров, делающих запросы
-H - строка, которая добавится к заголовку Content-Type

### В основном опирался на функционал инструмента ffuf
