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
        // "github.com/cydev/zero"
        "time"
        "sync"
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

var wg sync.WaitGroup

// func doPost(url string, payload []byte) ([]byte, error) {
//   req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
//   req.Header.Set("X-Custom-Header", "myvalue")
//   req.Header.Set("Content-Type", "application/json")
//
//   client := &http.Client{}
//   res, err := client.Do(req)
//   if err != nil {
//     panic (err)
//   }
//   defer res.Body.Close()
//   body, err := ioutil.ReadAll(res.Body)
//   return body, err
// }

func getBlockHeight() int{
  // get block height
  var payLoadBlockHeight = []byte(`{"method": "get_chain_block_height","jsonrpc": "2.0","id": 0}`)
  body, err := doPost("http://public.coindaddy.io:4100", payLoadBlockHeight)
  var blockHeight = new(resBlockHeight)
  err = json.NewDecoder(bytes.NewReader(body)).Decode(&blockHeight)
  if err != nil {
    panic(err.Error)
  }
  return blockHeight.Result
}

var total_txes_vacus int = 0
var total_txes int32 = 0
func getBlockInfo(blockIndexesArr []int, db *pg.DB) {
  //get Block info with specified block index
  blockIndexes := strings.Trim(strings.Replace(fmt.Sprint(blockIndexesArr), " ", ", ", -1), "[]")
  // fmt.Println(blockIndexes)
  var payLoadBlockInfo = []byte(`{"method": "get_blocks","params":{"block_indexes":[` + blockIndexes + `]},"jsonrpc": "2.0","id": 0}`)
  // var request_start = time.Now()
  body, err := doPost("http://rpc:1234@public.coindaddy.io:4000", payLoadBlockInfo)
  // body, err := doPost("http://rpc:rpc@35.200.18.50:14000/api/", payLoadBlockInfo)
  if err != nil {
    panic(err.Error)
  }
  // var request_end = time.Now()
  // fmt.Println("request time: " + request_end.Sub(request_start).String())
  var Blocks = new(BlockInfo)
  err = json.NewDecoder(bytes.NewReader(body)).Decode(&Blocks)
  if err!=nil {
    panic(err.Error)
  }
  var binding = new(Binding)
  // var count = 0
  for _, v := range Blocks.Result {
    for _, message := range v.Messages {
      wg.Add(1)
      go func(message Message, v Block){
        total_txes += 1
        time_tx1 := time.Now()
        if !strings.Contains(message.Bindings, "VACUS") || !strings.EqualFold(message.Category, "sends") {
          wg.Done()
          return
        }
        total_txes_vacus += 1
        err := json.NewDecoder(bytes.NewBufferString(message.Bindings)).Decode(&binding)
        if err != nil {
          panic(err)
        }
        binding.CreatedAt = v.BlockTime
        if !strings.EqualFold(binding.Status, "valid") {
          wg.Done()
          return
        }
        err = db.Insert(binding)
        if err != nil {
        }
        fmt.Println("time per tx: " + time.Now().Sub(time_tx1).String())
        wg.Done()
        // fmt.Println(binding)
        // fmt.Println(message.Bindings)
      }(message, v)
    }
  }
  wg.Done()
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

  // err := createSchema(db)
  // if err != nil {
  //   // panic(err)
  // }
  // list_address := []string{ "...", "..." }
  // ch := make(chan Block)
  // defer close(ch)

  var last_block_index = getBlockHeight()
  time_block_height := time.Now()
  block_height_execution := time_block_height.Sub(start)
  fmt.Println("block height cost: " + block_height_execution.String()) // ~ 2s
  // var block_indexes_1 = initArray(last_block_index, 249)
  wg.Add(50)
  for i:=0; i<50; i++{
    go getBlockInfo(initArray(last_block_index - i*40, 40), db)
  }
  wg.Wait()
  end_time := time.Now()
  total := end_time.Sub(start)
  fmt.Println("total Time execution: " + total.String())
  fmt.Print("total_txes: ")
  fmt.Println(total_txes)
  fmt.Println("total_txes of VACUS: " + strconv.Itoa(total_txes_vacus))
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

/* report:
env: test-Net
blocks: 1000
txes: 1206
txes_vacus: 5
time get block height: 0.677s
time_total: ~2s

env: test-Net
blocks: 2000
txes: 4252
txes_vacus: 29
time get block height: 0.68s
time: ~3.23

env: main-Net
blocks: 2000
txes: 35352
txes_vacus: 163
time to get block height: 0.556s
total time: 17.3s

env: main-Net
blocks: 1000
txes: 15703
txes_vacus: 73
time to get block height: 0.83s
total time: 7.67s
*/
