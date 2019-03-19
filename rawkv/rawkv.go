package rawkv

import (
	"flag"
	"fmt"
	"github.com/pingcap/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tikv/client-go/config"
	"github.com/tikv/client-go/raw"
	"github.com/yuin/gopher-lua"
	"strings"
)

var (
	pdAddr    = flag.String("pd", "0.0.0.0:2379", "pd address:localhost:37187")

)

func Loader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)
	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"newClient": newClient,
	"double": double,
	"put": put,
	"get": get,
	"assertEquals": assertEquals,
}

var clients = make(map[int]*raw.RawKVClient)

var nextId = 0

func double(L *lua.LState) int {
	lv := L.ToInt(1)             /* get argument */
	L.Push(lua.LNumber(lv * 2)) /* push result */
	return 1
}

func newClient(L *lua.LState) int {
	cli, err := raw.NewRawKVClient(strings.Split(*pdAddr, ","), config.Security{

	})
	if err != nil {
		log.Fatal(err)
	}
	currentId := nextId
	nextId += 1
	clients[currentId] = cli
	//fmt.Println(currentId)
	table := L.NewTable()
	fmt.Println(table)
	table.RawSetString("id", lua.LNumber(currentId))
	L.Push(table)
	return 1
}

func put(L *lua.LState) int {
	key := L.ToString(2)
	value := L.ToString(3)
	fmt.Println("key", key)
	fmt.Println("value", value)
	lv := L.Get(-3)

	if tbl, ok := lv.(*lua.LTable); ok {
		// lv is LTable
		currentId := tbl.RawGetString("id")
		if cId, ok := currentId.(lua.LNumber); ok {

			cli := clients[int(cId)]

			err := cli.Put([]byte(key), []byte(value))
			if err != nil {
				log.Fatal(errors.ErrorStack(err))
			}
		}
	}

	return 0
}

func get(L *lua.LState) int {
	lv := L.Get(-2)

	if tbl, ok := lv.(*lua.LTable); ok {
		// lv is LTable
		currentId := tbl.RawGetString("id")
		if cId, ok := currentId.(lua.LNumber); ok {

			cli := clients[int(cId)]
			key := L.ToString(2)
			res, err := cli.Get([]byte(key))
			if err != nil {
				log.Fatal(errors.ErrorStack(err))
			}
			L.Push(lua.LString(string(res)))
		}
	}

	return 1
}

func assertEquals(L *lua.LState) int {
	value1 := L.ToString(1)
	value2 := L.ToString(2)
	if value1 != value2 {
		panic("assertEqual fail")
	}

	return 0
}