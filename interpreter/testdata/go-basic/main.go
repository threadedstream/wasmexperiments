package main

import (
	"fmt"
	"log"
	"os"

	"github.com/threadedstream/wasmexperiments/api"
)

func main() {
	//path := flag.String("path", "", "path to a wasm binary")
	wapi, err := api.NewWasmApi(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	res, err := wapi.Call("fac", 5)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Result is %v", res)
}
