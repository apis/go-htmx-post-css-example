package wsNotifications

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type WsNotificationServer interface {
	Handler(responseWriter http.ResponseWriter, request *http.Request)
	Publish(message []byte)
}

const subscriberMessageBuffer = 16

type server struct {
	mutex       sync.Mutex
	subscribers map[*serverSubscriber]any
}

func NewServer() WsNotificationServer {
	wsNotificationServer := &server{
		subscribers: make(map[*serverSubscriber]any),
	}
	return wsNotificationServer
}

type serverSubscriber struct {
	messageChannel chan []byte
	closeSlow      func()
}

func (instance *server) Handler(responseWriter http.ResponseWriter, request *http.Request) {
	err := instance.subscribe(responseWriter, request)

	if errors.Is(err, context.Canceled) {
		return
	}

	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func (instance *server) subscribe(w http.ResponseWriter, r *http.Request) error {
	websocketConnection, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	defer func() {
		err := websocketConnection.CloseNow()
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	subscriber := &serverSubscriber{
		messageChannel: make(chan []byte, subscriberMessageBuffer),
		closeSlow: func() {
			if websocketConnection != nil {
				err := websocketConnection.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
				if err != nil {
					log.Printf("%v", err)
				}
			}
		},
	}

	instance.addSubscriber(subscriber)
	defer instance.deleteSubscriber(subscriber)

	ctx := websocketConnection.CloseRead(context.Background())

	for {
		select {
		case message := <-subscriber.messageChannel:
			err := writeTimeout(ctx, time.Second*5, websocketConnection, message)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (instance *server) Publish(message []byte) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	for subscriber := range instance.subscribers {
		select {
		case subscriber.messageChannel <- message:
		default:
			go subscriber.closeSlow()
		}
	}
}

func (instance *server) addSubscriber(subscriber *serverSubscriber) {
	instance.mutex.Lock()
	instance.subscribers[subscriber] = struct{}{}
	instance.mutex.Unlock()
}

func (instance *server) deleteSubscriber(subscriber *serverSubscriber) {
	instance.mutex.Lock()
	delete(instance.subscribers, subscriber)
	instance.mutex.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, websocketConnection *websocket.Conn, msg []byte) error {
	writeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return websocketConnection.Write(writeCtx, websocket.MessageText, msg)
}
