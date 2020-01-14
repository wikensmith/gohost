package main

import (
	"fmt"
	"gohost/queue"
)

func init() {
	Workers["YS.机票.国内.询价.wiken.DEBUG"] = func(context queue.Context) string {
		fmt.Println(string(context.QueueObj.Body))
		return "end"
	}
}
func main() {
	Start()
}
