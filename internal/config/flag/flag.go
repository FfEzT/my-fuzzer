package config

import (
	"flag"
	"fmt"
)

const (
  AND = iota
  OR
)

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
	Mode int
	// Regexp string
	// Lines  string
	// Size   string
	// Status string
	// Words  string
	// Time   ForTime
}

// type ForTime struct {
//   Operation int // use enum for GREATER or LESS
//   Number int
// }

func GetConfig() (result *Config, err error) {
  result = &Config{}

  target      := flag.String("u", "foo", "host target")
  worlistPath := flag.String("w", "", "Path to wordlist")
  payload     := flag.String("d", "", "payload")

  workerCount := flag.Int("c", 3, "Count of workers (default: 3)")
  contentType := flag.String("H", "text/plain", "Content-Type")
  mode        := flag.String("fmode", "and", "Filter set operator. Either of: and, or (default: and)")
  // regexp      := flag.String("fr", "", "Filter regexp")
  // time        := flag.String("ft", "***", "Filter by number of milliseconds to the first response byte, either greater or less than. EG: >100 or <100")
  // lines       := flag.String("fl", "***", "Filter HTTP status codes from response. Comma separated list of codes and ranges")
  // size        := flag.String("fs", "***", "Filter HTTP response size. Comma separated list of sizes and ranges")
  // status      := flag.String("fc", "***", "Filter HTTP status codes from response. Comma separated list of codes and ranges")
  // words       := flag.String("fw", "***", "Filter by amount of words in response. Comma separated list of word counts and ranges")

  flag.Parse()
  if *target == "" || *worlistPath == "" || *payload == "" {
    err = fmt.Errorf("usage: fuzzer -u host -w /path/to/wordlist -d \"user=admin&password=FUZZ\" [-fs 7,10-15 ; fmode or] ")
    result = nil

    return
  }

  result.Target      = *target
  result.WorlistPath = *worlistPath
  result.Payload     = *payload
  result.WorkerCount = *workerCount
  result.ContentType = *contentType
  // result.Filter.Regexp = *regexp

  if result.Filter.Mode, err = parseMod(*mode); err != nil {
    return
  }
  // if result.Filter.Time, err = parseTime(*time); err != nil {
  //   return
  // }
  // result.Filter.Lines    =
  // result.Filter.Size     =
  // result.Filter.Status   =
  // result.Filter.Words    =

  return
}

func parseMod(str string) (res int, err error) {
  if str == "or" {
    res = OR
  } else if str == "and" {
    res = AND
  } else {
    err = fmt.Errorf("[-fmode or/and]")
  }

  return
}