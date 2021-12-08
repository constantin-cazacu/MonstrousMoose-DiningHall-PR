package main

import (
	"fmt"
	"net/http"
	"os"
)
const WaiterNum = 1
const TableNum = 1


var kitchenHost = "http://localhost"
var tableList = newTableList()
var waiterList = newWaiterList()

func main() {

	args := os.Args

	if len(args) > 1{
		//Set the docker internal host
		kitchenHost = args[1]
	}

	fmt.Println("Dining hall is up and running!")

	for _, waiter := range waiterList {
		go waiter.work()
	}

	http.HandleFunc("/delivery",deliveryHandler)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}

