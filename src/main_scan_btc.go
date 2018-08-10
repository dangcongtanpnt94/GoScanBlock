package main

import (
        "fmt"
        "handlers"
        "models"
        // "strings"
        // "github.com/go-pg/pg"
        // "github.com/go-pg/pg/orm"
        // "strconv"
        // "github.com/cydev/zero"
        "time"
        "sync"
        "runtime"
)

var wg sync.WaitGroup

func main() {
  numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

  fmt.Println("=================")
  start := time.Now()

  var last_block_index = handlers.GetBlockHeight()

  block_indexes := handlers.InitArray(last_block_index, 498)


  // get all block_hashes
  var block_hashes []string
  for _, value := range block_indexes {
    wg.Add(1)
    go func(value int) {
      block_hashes = append(block_hashes, handlers.GetBlockHash(value))
      wg.Done()
    }(value)
  }
  wg.Wait()

  // get all block_hashes
  var blocks []models.Block
  for _, block_hash := range block_hashes {
    wg.Add(1)
    go func(block_hash string) {
      blocks = append(blocks, handlers.GetBlock(block_hash))
      wg.Done()
    }(block_hash)
  }
  wg.Wait()
  // fmt.Println(blocks)

  // scan all tx_es
  var txes []models.Tx
  for _, block := range blocks {
    // fmt.Println(block.Tx)
    for _, tx_hash := range block.Tx {
      wg.Add(1)
      go func(tx_hash string) {
        txes = append(txes, handlers.GetTx(tx_hash))
        wg.Done()
      }(tx_hash)
    }
  }
  wg.Wait()
  time_block_height := time.Now()
  block_height_execution := time_block_height.Sub(start)
  fmt.Println(fmt.Sprintf("Times cost %s \n last_block_index: %d", block_height_execution.String(), last_block_index)) // ~ 2s
  fmt.Printf("total txes: %d \n", len(txes))

}

// func createSchema(db *pg.DB) error {
//   for _, model := range []interface{}{ &Binding{}, }{
//     err := db.CreateTable(model, &orm.CreateTableOptions{
//       Temp: false,
//     })
//     if err != nil {
//       return err
//     }
//   }
//   return nil
// }

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

// ====================== try channel
  // block_hashes2 := make(chan string, 500)
  // for _, value := range block_indexes {
  //   wg.Add(1)
  //   go func(value int, result chan<- string ) {
  //     block_hashes2 <- handlers.GetBlockHash(value)
  //     wg.Done()
  //   }(value, block_hashes2)
  // }
  // wg.Wait()
  // close(block_hashes2)
  //
  // // for i:=1; i<= 500; i++ {
  // //   fmt.Println(<-block_hashes2)
  // // }
  //
  // blocks2 := make(chan models.Block, len(block_hashes2))
  // for block_hash := range block_hashes2 {
  //   wg.Add(1)
  //   go func(block_hash string, result chan<- models.Block) {
  //     blocks2 <- handlers.GetBlock(block_hash)
  //     wg.Done()
  //   }(block_hash, blocks2)
  // }
  // wg.Wait()
  // close(blocks2)
  // fmt.Printf("total blocks: %d \n", len(blocks2))
  // // for i:=1; i<=len(blocks2); i++ {
  // //   fmt.Println(<-blocks2)
  // // }
  //
  // txes2 := make(chan models.Tx, 10000)
  // for block := range blocks2 {
  //   // fmt.Println(block.Tx)
  //   for _, tx_hash := range block.Tx {
  //     wg.Add(1)
  //     go func(tx_hash string, result chan<- models.Tx) {
  //       txes2 <- handlers.GetTx(tx_hash)
  //       wg.Done()
  //     }(tx_hash, txes2)
  //   }
  // }
  // wg.Wait()
  // close(txes2)
// ====================== try channel
