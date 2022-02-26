package gateway

import (
	"bothoi/app_command"
	"bothoi/models"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

func dispatchHandler(c *websocket.Conn, payload models.GatewayPayload) {
	switch payload.T {
	case "READY":
		mapstructure.Decode(payload.D, &activeSessionState)
		sessionReady.Done()
	case "INTERACTION_CREATE":
		sessionReady.Wait()
		var data models.Interaction
		mapstructure.Decode(payload.D, &data)
		app_command.MapInteraction(&data)
	}
}
