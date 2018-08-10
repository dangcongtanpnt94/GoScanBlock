package models

type ResBlockHeightInfo struct {
  Result ResBlockHeight `json:"result"`
}

type ResBlockHeight struct {
  Blocks int `json:"blocks"`
}

type BlockHash struct {
  Result string `json:"result"`
}

type ResultBlock struct {
  Result Block `json:"result"`
}

type Block struct {
  Hash string `json:"hash"`
  Confirmations int `json:"confirmations"`
  Tx []string `json:"tx"`
}

type ResultTx struct {
  Result Tx `json:"result"`
}

type Tx struct {
  Txid string `json:"txid"`
  Vouts []Vout `json:"vout"`
}

type Vout struct {
  Value float32 `json:"value"`
  N int `json:"n"`
  Info VoutInfo `json:"scriptPubkey"`
}

type VoutInfo struct {
  addresses []string `json:"addresses"`
}
