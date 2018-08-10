package main

import (
  "fmt"
  "sync"
)

var wg sync.WaitGroup

func m(){
  fmt.Println("test1")
  wg.Done()
}

func n(){
  fmt.Println("test2")
  wg.Done()
}

func main(){
  wg.Add(10)
  for i:=0; i<=5; i++ {
    go m()
    go n()
  }
  wg.Wait()
}
