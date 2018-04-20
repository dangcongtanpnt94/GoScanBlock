package main

import (
        "fmt"
        "net/http"
        "bytes"
        "encoding/json"
        "io/ioutil"
        "strings"
        "github.com/go-pg/pg"
        "github.com/go-pg/pg/orm"
        "strconv"
        "github.com/cydev/zero"
        "time"
)

type Binding struct {
  Asset string `json:"asset"`
  Block_index int32 `json:"block_index"`
  Destination string `json:"destination"`
  Nemo string `json:"nemo"`
  Quantity int `json:"quantity"`
  Source string `json:"source"`
  Status string `json:"status"`
  TxHash string `json:"tx_hash"`
  TxIndex int32 `json:"tx_index" sql:",pk"`
  CreatedAt int64
}
type resBlockHeight struct {
  Id int `json:"id"`
  Jsonrpc string `json:"jsonrpc"`
  Result int `json:"result"`
}

type Message struct {
  MessageIndex int32 `json:"message_index"`
  Bindings string `json:"bindings"`
  Timestamp int64 `json:"timestamp"`
  Category string `json:"category"`
  BlockIndex int `json:"block_index"`
  Command string `json:"command"`
}

type Block struct {
  MessagesHash string `json:"nessages_hash"`
  Messages []Message `json:"_messages"`
  BlockHash string `json:"block_hash"`
  difficulty float64 `json:"difficulty"`
  BlockIndex int `json:"block_index"`
  BlockTime int64 `json:"block_time"`
  LedgerHash string `json:"ledger_hash"`
  TxlistHash string `json:"txlist_hash"`
  PreviousBlockHash string `json:"previous_block_hash"`
}

type BlockInfo struct{
  Id int `json:"id"`
  Jsonrpc string `json:"jsonrpc"`
  Result []Block `json:"result"`
}

func doPost(url string, payload []byte) ([]byte, error) {
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

func getBlockHeight() int{
  // get block height
  var payLoadBlockHeight = []byte(`{"method": "get_chain_block_height","jsonrpc": "2.0","id": 0}`)
  body, err := doPost("http://public.coindaddy.io:14100", payLoadBlockHeight)
  var blockHeight = new(resBlockHeight)
  err = json.NewDecoder(bytes.NewReader(body)).Decode(&blockHeight)
  if err != nil {
    panic(err.Error)
  }
  return blockHeight.Result
}

func getBlockInfo(blockIndexes int, c chan Block) {
  //get Block info with specified block index
  // blockIndexes := strings.Trim(strings.Replace(fmt.Sprint(blockIndexesArr), " ", ", ", -1), "[]")
  var payLoadBlockInfo = []byte(`{"method": "get_blocks","params":{"block_indexes":[` + strconv.Itoa(blockIndexes) + `]},"jsonrpc": "2.0","id": 0}`)
  body, err := doPost("http://rpc:1234@public.coindaddy.io:14000", payLoadBlockInfo)
  var Blocks = new(BlockInfo)
  err = json.NewDecoder(bytes.NewReader(body)).Decode(&Blocks)
  if err!=nil {
    panic(err.Error)
  }
  c <- Blocks.Result[0]
  // var blockHeight
}

// init an Array with value is a sequence number
// example: initArray(10, 5) => 10 9 8 7 6
func initArray(max, leng int) []int {
  a := make([]int, leng)
  for i := range a {
      a[i] = max - i
  }
  return a
}

func main() {
  // connect DB
  start := time.Now()
  db := pg.Connect(&pg.Options{
    User: "dangcongtan",
    Password: "",
    Database:"postgres",
  })
  defer db.Close()

  err := createSchema(db)
  if err != nil {
    // panic(err)
  }
  // list_address := []string{ "...", "..." }
  ch := make(chan Block)
  defer close(ch)

  var last_block_index = getBlockHeight()
  // var block_indexes_1 = initArray(last_block_index, 249)
  for i:=0; i<1000; i++{
    go getBlockInfo(last_block_index - i, ch)
  }
  // var blocks = getBlockInfo(block_indexes_1, ch)
  // fmt.Printf("%v", blocks.Result)
  var binding = new(Binding)
  var count = 0
  for {
    v := <-ch
    count++
    if count == 1000 {
      break
    }
    if zero.IsZero(v) {
      break
    }
    if len(v.Messages) == 0 {
      continue
    }
    for _, message := range v.Messages {
      if !strings.Contains(message.Bindings, "VACUS") {
        continue
      }
      if !strings.EqualFold(message.Category, "sends") {
        continue
      }
      err := json.NewDecoder(bytes.NewBufferString(message.Bindings)).Decode(&binding)
      if err != nil {
        panic(err)
      }
      binding.CreatedAt = v.BlockTime
      if !strings.EqualFold(binding.Status, "valid") {
        continue
      }
      err = db.Insert(binding)
      if err != nil {
      }
      fmt.Println(binding)
      // fmt.Println(message.Bindings)
    }
  }
  end_time := time.Now()
  total := end_time.Sub(start)
  fmt.Println("total Time execution: " + total.String())


  // var m = getBlockInfo([]int{518869, 518870})
  // fmt.Println(m.Result[0].Messages[0].Bindings)
  // getBlockHeight()
  // m :=4
  // var payLoadBlockInfo = `{"method": "get_block_info","params":{"block_index":` + strconv.Itoa(m) + `"jsonrpc": "2.0","id": 0}`
  // var m = []int{518869, 518870}
  // fmt.Println()
  // fmt.Println(payLoadBlockInfo)
  //
  // url = "https://xchain.io/api/asset/vacus"
  // fmt.Println("URL:>", url)
  //
  // var jsonStr = []byte(`{}`)
  // res, err := http.Get()
  // fmt.Println(res)
  // // fmt.Println(res.status)
  // if err != nil {
  //   panic(err.Error())
  // }
  // // fmt.Println(res.Body)
  // body, err := ioutil.ReadAll(res.Body)
  // if err != nil {
  //   panic(err.Error())
  // }
  // defer res.Body.Close()
}

func createSchema(db *pg.DB) error {
  for _, model := range []interface{}{ &Binding{}, }{
    err := db.CreateTable(model, &orm.CreateTableOptions{
      Temp: false,
    })
    if err != nil {
      return err
    }
  }
  return nil
}
