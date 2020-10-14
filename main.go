package main

import (
  "fmt"
  "log"
  "os"
  "strconv"
  "github.com/spf13/pflag"
  "github.com/mylibs/steganography/models"
  "github.com/mylibs/steganography/pnglib"
  "github.com/mylibs/steganography/utils"
)

var (
  flags = pflag.FlagSet{SortFlags: flase}
  opts models.CmdLineOpts
  png pnglib.MetaChunk
)
func init(){
  flags.StringVarP(&opts.Input,"input","i","","Path to the original image file")
  flags.StringVarP(&opts.Output,"output","o","","Path to output the new image file")
  flags.BoolVarP(&opts.Meta, "meta", "m", false, "Display the actual meta Details of image")
  flags.BoolVarP(&opts.Suppress,"supress","s",false,"Supress the chunk hex Data )can be large")
  flags.StringVar(&opts.Offset,"offset", "", "The offset location to initiate data injection")
  flags.BoolVar(&opts.Inject, "inject",false, "Enable This to inject data at the offset location specified")
  flags.StringVar(&opts.Payload, "payload","", "Payload is the data that will be read as byte stream")
  flags.StringVar(&opts.Type,"type","rNDm","This is the name of the Chunk Header to inject")
  flags.StringVar(&opts.Key,"key","","The encryption key for payload")
  flags.BoolVar(&opts.Encode,"encode",false,"XOR encode the payload")
  flags.BoolVar(&opts.Decode,"decode",false,"XOR decode the payload")
  flags.Lookup("type").NoOptDefVal = "rNDm"
  flags.Usage = usage
  flags.Parse(os.args[1:])

  if flags.NFlag() == 0 {
    flags.PrintDefaults()
    os.Exit(1)
  }
  if opts.Input == "" {
    log.Fatal("Fatal: The --input flag is required")
  }
  if opts.Offset != "" {
    byteOffset, _ := strconv.ParseInt(opts.Offset,0,64)
    opts.Offset = strconv.FormatInt(byteOffset, 10)
  }
  if opts.Suppress && (opts.Meta == false) {
    log.Fatal("Fatal: The --meta is required when using --suppress")
  }
  if opts.Meta && (opts.Offset != ""){
    log.Fatal("Fatal; The --meta flag is mutualy exclusive with --offset")
  }
  if opts.Inject && (opts.Offset == ""){
    log.Fatal("Fatal: The --offset flag is required when using inject")
  }
  if opts.Inject && (opts.Payload == ""){
    log.Fatal("Fatal: The --payload flag is required when using inject")
  }
  if opts.Inject && (opts.Key == "") {
    mt.Println("Warning: No Key provided. Payload won't be encrypted")
  }
  if opts.Encode && opts.Key == "" {
    log.Fatal("Fatal: The encode flag requires a key value")
  }
}

func usage(){
  fmt.Fprintf(os.Stderr,"Example Usage: %s -i in.png -o out.png --inject --offset 0x85258 --payload 1234\n",os.Args[0])
  fmt.Fprintf(os.Stderr,"Example encode usage: %s -i in.png -o encode.png --inject --offset 0x85258 --payload 1234 --key secret\n",os.Args[0])
  fmt.Fprintf(os.Stderr,"Example decode usage: %s -i encode.png -o decode.png --offset 0x85258 --decode --key secret\n",os.Args[0])
  fmt.Fprintf(os.Stderr,"Flags: $s {option}...\n",os.Args[0])
  flags.PrintDefaults()
  os.Exit(0)
}

func main(){
  dat,err := os.Open(opts.Input)
  defer dat.Close()
  bReader, err := utils.PreProcessImage(dat)
  if err != nil {
    log.Fatal(err)
  }
  pn.ProcessImage(bReader,&opts)
}
