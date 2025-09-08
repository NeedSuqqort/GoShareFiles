package counter

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

var cnt int
var mutex = &sync.Mutex{}

func health_check(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "OK")
}

func incrementCounter(writer http.ResponseWriter, request *http.Request) {
	mutex.Lock()
	cnt++
	fmt.Fprintf(writer, "%s", strconv.Itoa(cnt))
	mutex.Unlock()
}