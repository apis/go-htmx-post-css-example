package web

import (
	"htmx-example/internal/pkg/wsNotifications"
	"net/http"
)

type Handler struct {
	Request            func(request *http.Request, notificationServer wsNotifications.WsNotificationServer, simulatedDelay int) *Response
	SimulatedDelay     int
	NotificationServer wsNotifications.WsNotificationServer
}

func (instance Handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	response := instance.Request(request, instance.NotificationServer, instance.SimulatedDelay)
	response.Write(responseWriter)
}
