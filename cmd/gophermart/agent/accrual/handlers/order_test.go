package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	handlers "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/handlers/mocks"
	accrualStorage "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/storage"
	storage2 "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/mocks"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
)

var handlersOrderTestGlobalWg sync.WaitGroup

func TestOrdersProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	accrualStorage.NumberOfWorkers = 1

	tests := []struct {
		name              string
		returnFromStorage []user.OrderCheckAccrual
	}{
		{
			name: "test 1",
			returnFromStorage: []user.OrderCheckAccrual{{
				ID:        1,
				Number:    "1",
				Status:    "REGISTERED",
				Accrual:   100.1,
				UserLogin: "ramil",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlersOrderTestGlobalWg.Add(1)
			go func() {
				defer handlersOrderTestGlobalWg.Done()
				storageMock := storage2.NewMockStorager(ctrl)

				storageMock.EXPECT().
					GetAllOrdersInStatuses(context.Background(), gomock.Any()).
					Return(tt.returnFromStorage, nil)

				storageMock.EXPECT().
					UpdateOrderAccrual(context.Background(), gomock.Any()).
					Return(nil)

				clientMock := handlers.NewMockClienter(ctrl)

				var buf bytes.Buffer
				json.NewEncoder(&buf).Encode(tt.returnFromStorage[0])

				mockBody := io.NopCloser(&buf)
				defer mockBody.Close()

				var mockHeader http.Header
				clientMock.EXPECT().
					SendRequest("GET", "http://localhost:8081/api/orders/1", gomock.Any()).
					Return(200, mockBody, mockHeader, nil)

				go func() {
					go OrdersProcess(clientMock, storageMock)
				}()
				time.Sleep(1100 * time.Millisecond)
				log.Println("close without error")
			}()
		})
	}
}

func TestProcessAccrual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	accrualStorage.NumberOfWorkers = 1

	tests := []struct {
		name              string
		ticker            *time.Ticker
		returnFromStorage []user.OrderCheckAccrual
	}{
		{
			name:   "test 1",
			ticker: time.NewTicker(time.Second),
			returnFromStorage: []user.OrderCheckAccrual{{
				ID:        1,
				Number:    "1",
				Status:    "REGISTERED",
				Accrual:   100.1,
				UserLogin: "ramil",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlersOrderTestGlobalWg.Add(1)
			go func() {
				defer handlersOrderTestGlobalWg.Done()
				storageMock := storage2.NewMockStorager(ctrl)

				storageMock.EXPECT().
					GetAllOrdersInStatuses(context.Background(), gomock.Any()).
					Return(tt.returnFromStorage, nil)

				storageMock.EXPECT().
					UpdateOrderAccrual(context.Background(), gomock.Any()).
					Return(nil)

				clientMock := handlers.NewMockClienter(ctrl)

				var buf bytes.Buffer
				json.NewEncoder(&buf).Encode(tt.returnFromStorage[0])

				mockBody := io.NopCloser(&buf)
				defer mockBody.Close()

				var mockHeader http.Header
				clientMock.EXPECT().
					SendRequest("GET", "http://localhost:8081/api/orders/1", gomock.Any()).
					Return(200, mockBody, mockHeader, nil)

				go ProcessAccrual(clientMock, storageMock, tt.ticker)
				time.Sleep(1100 * time.Millisecond)
				log.Println("close without error")
			}()

		})
	}
}

func TestSyncAccrual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name              string
		returnFromStorage []user.OrderCheckAccrual
		url               string
		ordersCh          chan user.OrderCheckAccrual
		worker            int
		retryAfterChan    chan int
	}{
		{
			name: "test 1",
			returnFromStorage: []user.OrderCheckAccrual{{
				ID:        1,
				Number:    "1",
				Status:    "REGISTERED",
				Accrual:   100.1,
				UserLogin: "ramil",
			}},
			url:            "http://localhost:8081",
			worker:         1,
			retryAfterChan: make(chan int, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accrualStorage.NumberOfWorkers = 1
			storageMock := storage2.NewMockStorager(ctrl)

			storageMock.EXPECT().
				UpdateOrderAccrual(context.Background(), gomock.Any()).
				Return(nil)

			clientMock := handlers.NewMockClienter(ctrl)

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(tt.returnFromStorage[0])

			mockBody := io.NopCloser(&buf)
			defer mockBody.Close()

			var mockHeader http.Header
			clientMock.EXPECT().
				SendRequest("GET", "http://localhost:8081/api/orders/1", gomock.Any()).
				Return(200, mockBody, mockHeader, nil)

			ordersCh := make(chan user.OrderCheckAccrual, 1)
			defer close(ordersCh)
			ordersCh <- tt.returnFromStorage[0]

			go SyncAccrual(clientMock, tt.url, ordersCh, storageMock, tt.worker)
			<-time.After(1100 * time.Millisecond)
			log.Println("close without error")
			handlersOrderTestGlobalWg.Wait()
		})
	}
}
