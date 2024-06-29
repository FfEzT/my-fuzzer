package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// TODO for FilterOptions.Mode
// const (
//   AND = iota
//   OR
// )

type Config struct {
  Target      string
  WorlistPath string
  Payload     string
  WorkerCount int
  ContentType string
  Method      string
  Filter      FilterOptions
}

type FilterOptions struct {
  // TODO add mod for "and/or" between filters
	// Mode int
	Regexp string
	Status ValuesAndRanges
	Lines  ValuesAndRanges
	Words  ValuesAndRanges
	Size   ValuesAndRanges
	// Time   ForTime
}

type ValuesAndRanges struct {
  Values []int
  Ranges []Range
}

type Range struct {
  LeftValue, RightValue int
}

// type ForTime struct {
//   Operation int // use enum for GREATER or LESS or Undefined (не фильтровать по времени)
//   Number int
// }

func GetConfig() (result *Config, err error) {
  result = &Config{}

  target      := flag.String("u", "foo", "host target")
  worlistPath := flag.String("w", "", "Path to wordlist")
  payload     := flag.String("d", "", "payload")

  workerCount := flag.Int("c", 3, "Count of workers (default: 3)")
  contentType := flag.String("H", "text/plain", "Content-Type")
  // mode        := flag.String("fmode", "and", "Filter set operator. Either of: and, or (default: and)")
  regexp      := flag.String("fr", "", "Filter regexp")
  // time        := flag.String("ft", "", "Filter by number of milliseconds to the first response byte, either greater or less than. EG: >100 or <100")
  lines       := flag.String("fl", "", "Filter HTTP status codes from response. Comma separated list of codes and ranges")
  size        := flag.String("fs", "", "Filter HTTP response size. Comma separated list of sizes and ranges")
  status      := flag.String("fc", "", "Filter HTTP status codes from response. Comma separated list of codes and ranges")
  words       := flag.String("fw", "", "Filter by amount of words in response. Comma separated list of word counts and ranges")

  flag.Parse()
  if *target == "" || *worlistPath == "" || *payload == "" {
    err = fmt.Errorf("usage: fuzzer -u host -w /path/to/wordlist -d \"user=admin&password=FUZZ\" [-fs 7,10-15 ; fmode or] ")
    result = nil

    return
  }

  result.Target        = *target
  result.WorlistPath   = *worlistPath
  result.Payload       = *payload
  result.WorkerCount   = *workerCount
  result.ContentType   = *contentType
  result.Filter.Regexp = *regexp

  // if result.Filter.Mode, err = parseMod(*mode); err != nil {
  //   return
  // }
  // if result.Filter.Time, err = parseTime(*time); err != nil {
  //   return
  // }

  // NOTE в целом некритично, если не обработаю ошибки
  result.Filter.Status, _ = parseValueAndRanges(*status)
  result.Filter.Lines,  _ = parseValueAndRanges(*lines)
  result.Filter.Words,  _ = parseValueAndRanges(*words)
  result.Filter.Size,   _ = parseValueAndRanges(*size)

  return
}

func parseValueAndRanges(str string) (ValuesAndRanges, error) {
  var result ValuesAndRanges
  ranges := make([]Range, 0)
  values := make([]int,   0)

  parts := strings.Split(str, ",")
  for _, part := range parts {
    if strings.Contains(part, "-") {
      bounds := strings.Split(part, "-")

      start, err := strconv.Atoi(bounds[0])
      if err != nil {
        return result, err
      }
      end, err := strconv.Atoi(bounds[1])
      if err != nil {
        return result, err
      }

      if start >= end {
        return result, fmt.Errorf("invalid range: start is bigger than end")
      }

      // Добавляем интервал в слайс
      ranges = append(ranges, Range{start, end})
    } else {
      value, err := strconv.Atoi(part)
      if err != nil {
        return result, err
      }

      values = append(values, value)
    }
  }

  result.Ranges = ranges
  result.Values  = values

  return result, nil
}

// func parseMod(str string) (res int, err error) {
//   if str == "or" {
//     res = OR
//   } else if str == "and" {
//     res = AND
//   } else {
//     err = fmt.Errorf("[-fmode or/and]")
//   }

//   return
// }