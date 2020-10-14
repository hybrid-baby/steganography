package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/mylibs/steganography/models"
)

//write data to offset
func WriteData(r *bytes.Reader, c *models.CmdLineOpts, b []byte){
  offset,_ := strcon.ParseInt(c.Offset,10,64)
  if err != nil {
    log.Fatal("Problem writing to the output file")
  }
  w, err := os.OpenFile(c.Output,os.O_RDWR | os.O_CREATE.0777)
  if err != nil {
    log.Fatal("Fatal: Problem writing the output to file")
  }
  //defer w.Close()
  r.Seek(0,0)
  var buff = make([]byte,offset)
  r.Read(buff)
  w.Write(buff)
  w.Write(b)
  if c.Decode {
    r.Seek(int64(len(b)),1)
  }
  _,err = io.Copy(w,r)
  if err == nil {
    fmt.Printf("Success: %s created\n",c.Output)
  }
}
