package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"log"
	"os"
)

//go:embed testdata/dummywasm/target/wasm32-unknown-unknown/debug/dummywasm.wasm
var greetWasm []byte

func main() {
	ctx := context.Background()

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	_, err := r.NewHostModuleBuilder("env").
		ExportFunction("log", logString).
		Instantiate(ctx, r)
	if err != nil {
		log.Panicln(err)
	}

	mod, err := r.InstantiateModuleFromBinary(ctx, greetWasm)
	if err != nil {
		log.Panic(err)
	}

	greet := mod.ExportedFunction("greet")
	greeting := mod.ExportedFunction("greeting")
	allocate := mod.ExportedFunction("allocate")
	deallocate := mod.ExportedFunction("deallocate")
	sendVec := mod.ExportedFunction("send_vec")

	name := os.Args[1]
	nameSize := uint64(len(name))

	results, err := allocate.Call(ctx, nameSize)
	if err != nil {
		log.Panic(err)
	}

	namePtr := results[0]

	defer deallocate.Call(ctx, namePtr, nameSize)

	if !mod.Memory().Write(ctx, uint32(namePtr), []byte(name)) {
		log.Panicf("Memory.Write(%d, %d) out of range of memory size %d",
			namePtr, nameSize, mod.Memory().Size(ctx))
	}

	_, err = greet.Call(ctx, namePtr, nameSize)
	if err != nil {
		log.Panic(err)
	}

	ptrSize, err := greeting.Call(ctx, namePtr, nameSize)
	if err != nil {
		log.Panic(err)
	}
	greetingPtr := uint32(ptrSize[0] >> 32)
	greetingSize := uint32(ptrSize[0])

	defer deallocate.Call(ctx, uint64(greetingPtr), uint64(greetingSize))

	if bytes, ok := mod.Memory().Read(ctx, greetingPtr, greetingSize); !ok {
		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
			greetingPtr, greetingSize, mod.Memory().Size(ctx))
	} else {
		fmt.Println("go >>", string(bytes))
	}

	rawPtr, err := sendVec.Call(ctx)
	if err != nil {
		log.Panic(err)
	}
	vecPtr := uint32(rawPtr[0] >> 32)
	vecSize := uint32(rawPtr[0])
	println(vecSize)
	fmt.Printf("0x%x", vecPtr)

	defer deallocate.Call(ctx, uint64(vecPtr))

	if bytes, ok := mod.Memory().Read(ctx, vecPtr, vecSize); !ok {
		log.Panic("unsuccessful memory read")
	} else {
		println(bytes)
	}
}

func logString(ctx context.Context, m api.Module, offset, byteCount uint32) {
	buf, ok := m.Memory().Read(ctx, offset, byteCount)
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", offset, byteCount)
	}
	fmt.Println(string(buf))
}
