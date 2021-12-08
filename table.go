package main

import (
	"math/rand"
	"sync"
	"time"
)

var tableId = 0

type Table struct {
	id           int
	readyToOrder bool
}

func newTable() Table {
	var ret Table
	ret.id = tableId
	tableId += 1
	ret.readyToOrder = true
	return ret
}

type TableList struct {
	tableArr []Table
	mx       sync.Mutex
}

func newTableList()TableList{
	var tableArr []Table
	for i := 0; i < TableNum; i++ {
		tableArr = append(tableArr, newTable())
	}
	var ret TableList
	ret.tableArr = tableArr
	return ret
}

var orderId = 0

func (t Table) generateOrder(waiter Waiter) *Order {
	ret := new(Order)
	ret.OrderId = orderId
	orderId += 1
	ret.TableId = t.id
	ret.Items = getItems()
	ret.Priority = rand.Intn(5)

	itemMaxPrep := 0
	for _, item := range ret.Items {
		if menu[item].prepTime > itemMaxPrep {
			itemMaxPrep = menu[item].prepTime
		}
	}
	ret.MaxWait = int(float32(itemMaxPrep) * 1.3)
	ret.PickUpTime = time.Now().Unix()
	return ret
}

func getItems() []int {
	var ret []int

	var itemNr = rand.Intn(10) + 1
	for i := 0; i < itemNr; i++ {
		ret = append(ret, rand.Intn(10)+1)
	}
	return ret
}
