package main

import (
	"fmt"
	"github.com/threadedstream/wasmexperiments/api"
	"log"
	"os"
)

func main() {
	//path := flag.String("path", "", "path to a wasm binary")
	wapi, err := api.NewWasmApi(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	res, err := wapi.Call("add", 10, 20)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Result is %v", res)
}
