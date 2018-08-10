package handlers

import (
  "fmt"
  "bytes"
  "encoding/json"
  "models"

)

func BitcoinCashRPC(payload []byte) []byte {
  body, err := DoPost("http://bitcoin-rpc:bitcoin-rpc-pass@35.200.108.6:18332", payload)
  if err != nil {
    panic(err.Error)
  }
  // fmt.Println(string(body[:])) // convert []byte to string
  return body
}

func GetBlockHeight() int{
  var payload = []byte(`{"method": "getinfo","jsonrpc": "2.0","id": 0}`)
  body := BitcoinCashRPC(payload)
  // fmt.Println(string(body[:])) // convert []byte to string
  var info = new(models.ResBlockHeightInfo)
  err := json.NewDecoder(bytes.NewReader(body)).Decode(&info)
  if err != nil {
    panic(err.Error)
  }
  return info.Result.Blocks
}

func GetBlockHash(block_index int) string {
  var payload = []byte(fmt.Sprintf(`{"method": "getblockhash","params": [%d], "jsonrpc": "2.0","id": 0}`, block_index))
  body := BitcoinCashRPC(payload)
  var blockHash = new(models.BlockHash)
  err := json.NewDecoder(bytes.NewReader(body)).Decode(&blockHash)
  if err != nil {
    panic(err.Error)
  }
  return blockHash.Result
}

func GetBlock(block_hash string) models.Block {
  var payload = []byte(fmt.Sprintf(`{"method": "getblock","params": ["%s"], "jsonrpc": "2.0","id": 0}`, block_hash))
  body := BitcoinCashRPC(payload)
  resultBlock := new(models.ResultBlock)
  err := json.NewDecoder(bytes.NewReader(body)).Decode(&resultBlock)
  if err != nil {
    panic(err.Error)
  }
  return resultBlock.Result
}

func GetTx(tx_hash string) models.Tx {
  var payload = []byte(fmt.Sprintf(`{"method": "getrawtransaction","params": ["%s", 1], "jsonrpc": "2.0","id": 0}`, tx_hash))
  body := BitcoinCashRPC(payload)
  resultTx := new(models.ResultTx)
  err := json.NewDecoder(bytes.NewReader(body)).Decode(&resultTx)
  if err != nil {
    panic(err.Error)
  }
  return resultTx.Result
}
