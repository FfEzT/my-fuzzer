package request

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Response struct {
  Payload string
  Status string
  Body string
  Duration time.Duration
  Err error
}

type Request struct {
  Method string
  Target string
  ContentType string
}

func SendRequest(payload string, request *Request) *Response {
  req, err := http.NewRequest(
    request.Method,
    request.Target,
    bytes.NewBuffer(
      []byte(payload),
    ),
  )
  if err != nil {
    return &Response{Payload: payload, Err:err}
  }
  req.Header.Set("Content-Type", request.ContentType)

  startTime := time.Now()

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return &Response{Payload: payload, Err:err}
  }
  defer resp.Body.Close()

  endTime := time.Now()


  body, _ := io.ReadAll(resp.Body)
  return &Response{
    payload,
    resp.Status,
    string(body),
    endTime.Sub(startTime),
    nil,
  }
}