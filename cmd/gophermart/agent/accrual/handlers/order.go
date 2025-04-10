package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	accrualStorage "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/flags"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	orderStatus "github.com/ramil063/firstgodiplom/internal/constants/status"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

func OrdersProcess(c Clienter, s storage.Storager) {

	duration := time.Duration(accrualStorage.OrderCheckTickerTimeInterval)
	ticker := time.NewTicker(duration * time.Second)

	go ProcessAccrual(c, s, ticker)
}

func ProcessAccrual(c Clienter, s storage.Storager, ticker *time.Ticker) {
	defer ticker.Stop()
	retryAfterChan := make(chan int, 1)
	var retryAfter int

	for {
		select {
		case retryAfter = <-retryAfterChan:
			time.Sleep(time.Duration(retryAfter) * time.Second)
			log.Printf("retryAfterChan Resuming work...\n")
			continue // Возвращаемся к отправке запросов
		case <-ticker.C:
			ordersCh := make(chan user.OrderCheckAccrual)
			go func() {
				defer close(ordersCh)
				orders, err := s.GetAllOrdersInStatuses([]int{orderStatus.NewID, orderStatus.ProcessedID})
				if err != nil {
					log.Println(err.Error())
				} else {
					for _, order := range orders {
						ordersCh <- order
					}
				}
			}()

			for worker := 0; worker < accrualStorage.NumberOfWorkers; worker++ {
				go SyncAccrual(c, flags.AccrualSystemAddress, ordersCh, s, worker, retryAfterChan)
			}
		}
	}
}

func SyncAccrual(c Clienter, url string, ordersCh chan user.OrderCheckAccrual, s storage.Storager, worker int, retryAfterChan chan int) {
	ordersCheckUrl := url + "/api/orders/"

	for order := range ordersCh {
		responseCode, body, header, err := c.SendRequest("GET", ordersCheckUrl+order.Number, []byte{})

		if err != nil {
			logger.WriteErrorLog("error while check order in accrual system: " + err.Error())
			log.Println("error while check order in accrual system" + err.Error())
		}

		log.Println(worker, "-worker; accrual-", responseCode, ordersCheckUrl+order.Number)

		if responseCode == http.StatusOK {
			var order accrualStorage.Order
			dec := json.NewDecoder(body)
			err := dec.Decode(&order)
			if err == nil {
				logMsg, _ := json.Marshal(order)
				logger.WriteInfoLog(string(logMsg))
				log.Println("accrual: 200, ", order.Order)
				err = s.UpdateOrderAccrual(order)
			}

			if err != nil {
				logger.WriteErrorLog(err.Error())
				log.Println("accrual: 200, decode error:", err.Error())
			}
		}

		if responseCode == http.StatusNoContent {
			log.Println("order not registered in accrual system")
			err = s.UpdateOrderCheckAccrualAfter(order.Number)
			if err != nil {
				logger.WriteDebugLog(err.Error())
			}
		}

		if responseCode == http.StatusTooManyRequests && header != nil {
			retryAfter := header.Get("Retry-After")
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					log.Printf("Worker %d: Received 429 status, will retry after %d seconds\n", worker, seconds)
					retryAfterChan <- seconds
				} else {
					log.Printf("Worker %d: Received 429 status, will retry after %d seconds\n", worker, accrualStorage.DefaultRetryAfterInterval)
					retryAfterChan <- accrualStorage.DefaultRetryAfterInterval
				}
			} else {
				log.Printf("Worker %d: Received 429 status, will retry after %d seconds\n", worker, accrualStorage.DefaultRetryAfterInterval)
				retryAfterChan <- accrualStorage.DefaultRetryAfterInterval
			}
		}

		if responseCode == http.StatusInternalServerError {
			log.Println("accrual: 500", body)
			logger.WriteErrorLog("error while check order in accrual system (500)")
		}

		body.Close()
	}
}
