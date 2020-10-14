package pnglib

import (
  "encoding/binary"
  "log"
  "fmt"
  "strconv"
  "bytes"
  "strings"
  "hash/crc32"
  "encoding/hex"
  "github.com/mylibs/steganography/utils"
  "github.com/mylibs/steganography/models"
)

//header data store
type Header struct {
  Header uint64
}

//CRC data chunk reresentation
type Chunk struct {
  Size uint32
  Type uint32
  Data []byte
  CRC unit32
}
type MetaChunk struct {
  Chk Chunk
  Offset   int64
}
//validate png by reading the header data
func (mc *MetaChunk) validate(b *bytes.Reader) {
  var header Header

  if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
    log.Fatal(err)
  }

  bArr := make([]byte,8)
  binary.BigEndian.PutUnit64(bArr,header.Header)
  if string(bArr[1:4]) != "PNG" {
    log.Fatal("Provided File isn't a PNG format")
  }else{
    fmt.Println("Valid PNG")
  }
}

func (mc *MetaChunk) readChunk(b *bytes.Reader){
  mc.readChunkSize(b)
  mc.readChunkType(b)
  mc.readChunkBytes(b, mc.Chk.Size)
  mc.readChunkCRC(b)
}

func (mc *MetaChunk) readChunkSize(b *bytes.Reader){
  if err := binary.Read(b,binary.BigEndian, &mc.Chk.Size)err != nil {
    log.Fatal(err)
  }
}

func (mc *MetaChunk) readChunkType(b *bytes.Reader){
  if err := binary.Read(b, binary.BigEndian, &mc.Chk.Type); err != nil {
    log.Fatal(err)
  }
}

func (mc *MetaChunk) readChunkBytes(b *bytes.Reader, cLen uint32) {
  mc.Chk.Data = make([]byte, cLen)
  if err := binary.Read(b,binary.BigEndian, &mc.Chk.Data); err != nil {
    log.Fatal(err)
    }
}

func (mc *MetaChunk) readChunkCRC(b *bytesReader) {
  if er := binary.Read(b, binary.BigEndian, &mc.Chk.CRC); err != nil {
    log.Fatal(err)
  }
}

func (mc *MetaChunk) getOffset(b *bytes.Reader) {
  offset, _ := b.Seek(0,1)
  mc.Offset = offset
}

func (mc *MetaChunk) chunkTypeToString() string {
  h := fmt.Sprintf("%x",mc.Chk.Type)
  decoded, _ := hex.DecodeString(h)
  result := fmt.Sprintf("%s",decoded)
  return result
}

func (mc *MetaChunk) checkCritType() string {
  fChar := string([]rune(mc.chunkTypeToString()),[0])
  if fChar ==  strings.ToUpper(fChar) {
    return "Critical"
  }
  return "Anciliary"
}

func createChunkSize() uint32{
  return uint32(len(mc.Chk.Data))
}

func (mc *MetaChunk) createChunkCRC() uint32 {
  bytesMSB := new(bytes.Buffer)
  if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type);err != nil {
    log.Fatal(err)
  }
  if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
    log.Fatal(err)
  }
  return crc32.ChecksumIEEE(bytesMSB.Bytes())
}

func (mc *MetaChunk) strToInt(s string) uint32{
  t := []byte(s)
  return binary.BigEndian.Uint32(t)
}

func (mc *MetaChunk) marshalData() *bytes.Buffer{
  bytesMSB := new(bytes.Buffer)
  if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Size); err != nil {
    log.Fatal(err)
  }
  if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
    log.Fatal(err)
  }
  if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil{
    log.Fatal(err)
  }
  if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.CRC); err != nil {
    log.Fatal(err)
  }
  return bytesMSB
}


//process Image
func (mc *MetaChunk) ProcessImage(b * bytes.Reader, c *models.CmdLineOpts) {
  mc.Validate(b)
  if(c.Offset != "") && (c.Encode == false && c.Decode == false){
    var m MetaChunk
    m.Chk.Data = []byte(c.Payload)
    m.Chk.Type = m.strToInt(c.Type)
    m.Chk.Size = m.createChunkSize()
    m.Chk.CRC = m.createChunkCRC()
    bm := m.marshalData()
    bmb := bm.Bytes()
    fmt.Printf("Payload Original?: % X\n",[]byte(c.Payload))
    fmt.Printf("Payload: %X\n",m.CHK.Data)
    utils.WriteData(b, c, bmb)
  }
  if (mc.Offset != "") && c.Encode {
    var m MetaChunk
    m.Chk.Data = utils.XorEncode([]byte(c.Payload), c.Key)
    m.Chk.Type = m.strToInt(c.Type)
    m.Chk.Size = m.createChunkSize()
    m.Chk.CRC = m.createChunkCRC()
    bm := m.marshalData()
    bmb := bm.Bytes()
    fmt.Printf("Payload Original: % X\n",[]byte(c.Payload))
    fmt.Printf("Payload encode: % X\n",m.Chk.Data)
    utils.WriteData(b, c, bmb)
  }
  if (c.Offset != "") && c.Decode {
    var m MetaChunk
    offset, _ := strconv,ParseInt(c.Offset, 10,64)
    b.Seek(offset, 0)
    m.readChunk(b)
    origData := m.Chk.Data
    m.Chk.Data = utils.XorDecode(m.Chk.Data,c.Key)
    m.Chk.CRC = m.createChunkCRC()
    bm := m.marshalData()
    bmb := bm.Bytes()
    fmt.Printf("Payload Original: % X\n",origData)
    fmt.Printf("Payload encode: % X\n",m.Chk.Data)
    utils.WriteData(b, c, bmb)
  }
  if c.Meta {
    count := 1
    var chunkType string
    for chunkType != endChunkType {
      mc.getOffset(b)
      mc.readChunk(b)
      fmt.Println("---- Chunk  #"+strconv.Itoa(count) + "----")
      fmt.Printf("Chunk offset: %#02x\n",mc.Offset)
      fmt.Printf("Chunk length: %s bytes\n",strconv.Itoa(int(mc.Chk.Size)))
      fmt.Printf("Chunk Type: %s\n",mc.chunkTypeToString())
      fmt.Printf("Chunk importance: %s\n",mc.checkCritType())
      if c.Suppress == false {
        fmt.Printf("Chunk Data: %#x\n",mc.Chk.Data)
      }else if c.Supress {
        fmt.Printf("Chunk Data: %s\n","suppressed")
      }
      fmt.Printf("Chunk CRC: %x\n",mc.Chk.CRC)
      chunkType = chunkTypeToString()
      count++
    }
  }
}
