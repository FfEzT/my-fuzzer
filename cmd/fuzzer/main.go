package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	config "fuzzer/internal/config/flag"
	"fuzzer/pkg/http/request"
)

func main() {
  var conf *config.Config
  var err error

  if conf, err = config.GetConfig(); err != nil {
    fmt.Println("Failed to get config:", err)
		return
  }

  file, err := os.Open(conf.WorlistPath)
  if err != nil {
    fmt.Println("Failed to open file:", err)
    return
  }
  defer file.Close()

  // Идея многопоточности
  // запускаем N воркеров и даем им каналы на получение аргументов и запись результатов запроса
  // запускаем горутину (назовем ее A), которая будет считывать с файла и в канал передавать аргументы, когда она закончит передавать, закроет канал. И будет ждать, пока не закончат работу воркеры (с помощью WaitGroup)
  // Все это время главный поток получает аргументы с воркеров
  // Воркеры закончат свою работу, когда A закроет канал с аргументами, после горутина A закроет канал с результатами запросов, тем самым освобождает main поток
  wg := sync.WaitGroup{}
  wg.Add(conf.WorkerCount)
  // сюда будет передаваться строка, которая будет отправляться в теле POST запроса
  arg := make(chan string, conf.WorkerCount)
  result := make(chan *request.Response, conf.WorkerCount)

  request := request.Request{
    ContentType: conf.ContentType,
    Target:      conf.Target,
    Method:      "POST",
  }

  for i := 0; i < conf.WorkerCount; i++ {
    go worker(&request, &wg, arg, result)
  }

  // TODO GraceFull ShutDown для чтения файлов и воркеров
  // done := make(chan os.Signal, 1)
  // signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

  // горутина, читающая файл
  go func() {
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      arg <- strings.ReplaceAll(
        conf.Payload,
        "FUZZ",
        scanner.Text(),
      )
    }
    close(arg)
    wg.Wait()
    close(result)
  }()

  startTime := time.Now()
  for response := range result {
    showResponse(response, conf)
  }
  endTime := time.Now()
  fmt.Println("\n\nTime of work", endTime.Sub(startTime))

}

func worker(req *request.Request, wg *sync.WaitGroup, payload <-chan string, result chan <- *request.Response) {
  // defer wg.Done()
  for word := range payload {
    result <- request.SendRequest(
      word,
      req,
    )
  }
  wg.Done()
}

func showResponse(resp *request.Response, conf *config.Config) (err error) {
  if resp.Err != nil {
    fmt.Println(resp.Payload, "\t", resp.Err)
    err = resp.Err
    return
  }

  size := len(resp.Body)
  words := len(strings.Fields(resp.Body))
  lines := len(strings.Split(resp.Body, "\n"))

  
  // наверное, выдавать надо, если нет указанного слова
  // тогда придется добавлять проверку на то, что аргумент == ""
  // TODO filter by RegExp
  if !strings.Contains(resp.Body, conf.Filter.Regexp) {
    return
  }

  // filter (фильтры убирают запросы из вывода)
  if isInValuesAndRanges(size, conf.Filter.Size) ||
     isInValuesAndRanges(lines, conf.Filter.Lines) ||
     isInValuesAndRanges(resp.Status, conf.Filter.Status) ||
     isInValuesAndRanges(words, conf.Filter.Words) {
    return
  }

  // TODO pretty output
  // tabwriter.Writer
  fmt.Printf(
    "%s\t[Status: %d, Size: %d, Words: %d, Lines: %d, Duration: %dms]\n",
    resp.Payload,
    resp.Status,
    size,
    words,
    lines,
    resp.Duration.Milliseconds(),
  )

  return
}

// TODO наверное стоит перенести в package config
// TODO извините за название
func isInValuesAndRanges(testable  int, filter config.ValuesAndRanges) bool {
  for _, value := range filter.Values {
    if testable == value {
      return true
    }
  }
  for _, rang := range filter.Ranges {
    if rang.LeftValue <= testable && testable <= rang.RightValue {
      return true
    }
  }
  return false
}