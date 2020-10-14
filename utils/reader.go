package utils

import (
  "bufio"
  "os"
  "bytes"
  "log"
)

//preprocess image reads a file image handle and return the buffer bytes
func PreProcessImage(datum *os.File)(*bytes.Reader, error) {
  stats,err := datum.Stat()
  if err != nil {
    return nil,err
  }
  var size = stats.Size()
  b := make([]byte,size)
  bufR := bufio.NewReader(datum)
  if err := bufR.Read(b); err != nil {
    log.Fatal(err)
  }
  bReader := bytes.NewReader(b)
  return bReader,err
}
