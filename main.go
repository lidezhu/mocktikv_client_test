package main

import (
	"mocktikv_client_test/rawkv"
	"path/filepath"

	//"github.com/pingcap/tidb/config"
	"github.com/yuin/gopher-lua"
	"log"
)




func main() {

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("rawkv", rawkv.Loader)
	dir, err := filepath.Abs(filepath.Dir("hello.lua"))
	if err != nil {
		log.Fatal(err)
	}


	if err := L.DoFile(dir+"/src/mocktikv_client_test/hello.lua"); err != nil {
		panic(err)
	}
}
