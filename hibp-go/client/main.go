package main

import (
  "bufio"
  "crypto/sha1"
  "fmt"
  "io"
  "log"
  "net/http"
  "strconv"
  "strings"
  "time"
)

func main() {

  const password = "123456"
  h := sha1.New()
  _, err := io.WriteString(h, password)
  if err != nil {
    log.Panicf("sha1 write string failed %v\n", err)
  }

  sha1hash := fmt.Sprintf("%x", h.Sum(nil))
  sha1hash = strings.ToUpper(sha1hash)
  httpClient := &http.Client{
    Timeout: time.Second * 10,
  }

  sha1hashFirst5 := sha1hash[0:5]
  sha1hashRest := sha1hash[5:]

  const url = "https://api.pwnedpasswords.com/range/"
  req, err := http.NewRequest("GET", url+sha1hashFirst5, nil)
  if err != nil {
    log.Panicf("creating request failed %v\n", err)
  }
  req.Header.Set("Add-Padding", "true")

  response, err := httpClient.Do(req)
  if err != nil {
    log.Panicf("http client get failed %v\n", err)
  }
  defer response.Body.Close()

  scanner := bufio.NewScanner(response.Body)
  found := false

  for scanner.Scan() {
    line := scanner.Text()
    appearances, err := strconv.Atoi(line[36:])

    if err != nil {
      log.Println(line)
      log.Panicf("conversion failed %s %v\n", line[37:], err)
    }

    if appearances == 0 {
      break
    }

    if line[:35] == sha1hashRest {
      fmt.Println("found it: ", appearances)
      found = true
      break
    }
  }

  if !found {
    fmt.Println("not found")
  }

}
