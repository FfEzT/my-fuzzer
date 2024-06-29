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
    fmt.Println("Failed to open file :", err)
    return
  }
  defer file.Close()

  // Идея многопоточности
  // запускаем N воркеров и даем им каналы на получение аргументов и запись результатов запроса
  // запускаем горутину (назовем ее A), которая будет считывать с файла и в канал передавать аргументы, когда она закончит передавать, закроет канал. И будет ждать, пока не закончат работу воркеры (с помощью WaitGroup)
  // Все это время главный поток получает аргументы с воркеров
  // Воркеры закончат свою работу, горутина A закроет канал с результатами запросов, тем самым освобождает main поток
  wg := sync.WaitGroup{}
  wg.Add(conf.WorkerCount)
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
  for i := range result {
    // TODO prettier
    // TODO filter
    // fmt.Print(i.Payload, "\t\t\t")
    if i.Err != nil {
      fmt.Println(i.Err)
      continue
    }

    // size := len(i.Body)
    // words := len(strings.Fields(i.Body))
    // lines := len(strings.Split(i.Body, "\n"))

    // fmt.Printf(
    //   "[Status: %s, Size: %d, Words: %d, Lines: %d, Duration: %dms]\n",
    //   i.Status,
    //   size,
    //   words,
    //   lines,
    //   i.Duration.Milliseconds(),
    // )
  }
  fmt.Println("Time of work", time.Now().Sub(startTime))

  // TODO Gracefull shutdown
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