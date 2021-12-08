package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var waiterId = 0

type Waiter struct {
	id int
}

func newWaiterList() []Waiter {
	var ret []Waiter
	for i := 0; i < WaiterNum; i++ {
		ret = append(ret, Waiter{waiterId})
		waiterId += 1
	}
	return ret
}

func (w Waiter) work() {
	for {
		tableList.mx.Lock()
		for i, _ := range tableList.tableArr {
			table := tableList.tableArr[i]
			if table.readyToOrder == true {
				tableList.tableArr[i].readyToOrder = false
				order := table.generateOrder(w)
				sendOrder(order)
				fmt.Println("Order", order.OrderId, "was sent from table:",table.id)

				break
			}
		}
		tableList.mx.Unlock()

		select {
		case delivery := <-deliveryChan:
			//Serve delivery to the required table
			tableList.mx.Lock()
			tableList.tableArr[delivery.TableId].readyToOrder = true
			now := time.Now().Unix()

			rating := 0
			maxWaitF := float64(delivery.MaxWait)
			timeWaitedF := float64(now - delivery.PickUpTime)
			if maxWaitF > timeWaitedF {
				rating += 1
			}
			if maxWaitF*1.1 > timeWaitedF {
				rating += 1
			}
			if maxWaitF*1.2 > timeWaitedF {
				rating += 1
			}
			if maxWaitF*1.3 > timeWaitedF {
				rating += 1
			}
			if maxWaitF*1.4 > timeWaitedF {
				rating += 1
			}

			fmt.Println("Order", delivery.OrderId, " | Rating: ", rating)
			tableList.mx.Unlock()
		default:
			break
		}
		time.Sleep(time.Second)

	}
}

func sendOrder(orderTicket *Order) bool {
	requestBody, marshallErr := json.Marshal(orderTicket)
	if marshallErr != nil {
		log.Fatalln(marshallErr)
	}

	request, newRequestError := http.NewRequest(http.MethodPost, kitchenHost+":8000"+"/order", bytes.NewBuffer(requestBody))
	if newRequestError != nil {
		fmt.Println("Could not create new request. Error:", newRequestError)
		log.Fatal(newRequestError)
	} else {
		response, doError := http.DefaultClient.Do(request)
		fmt.Println("Sending order to kitchen attempt")
		if doError != nil {
			fmt.Println("ERROR Sending request. ERR:", doError)
			log.Fatal(doError)
		}
		var responseBody = make([]byte, response.ContentLength)
		response.Body.Read(responseBody)
		fmt.Println("Response: ", string(responseBody))
		if string(responseBody) != "OK" {
			return false
		}
		return true
	}
	return true
}
