package handlers

import (
  "net/http"
  "bytes"
  "io/ioutil"
)

func DoPost(url string, payload []byte) ([]byte, error) {
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
  req.Header.Set("X-Custom-Header", "myvalue")
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    panic (err)
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return body, err
}

func DoGet(url string) ([]byte, error) {
  res, err := http.Get(url)
  if err != nil {
    panic (err)
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return body, err
}

// init an Array with value is a sequence number
// example: InitArray(10, 5) => 10 9 8 7 6
func InitArray(max, leng int) []int {
  a := make([]int, leng)
  for i := range a {
      a[i] = max - i
  }
  return a
}
