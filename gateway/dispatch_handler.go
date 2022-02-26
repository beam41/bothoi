package gateway

import (
	"bothoi/app_command"
	"bothoi/models"
	"bothoi/states"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

func dispatchHandler(c *websocket.Conn, payload models.GatewayPayload) {
	switch payload.T {
	case "READY":
		mapstructure.Decode(payload.D, &states.SessionState)
		states.SessionStateReady.Done()
	case "INTERACTION_CREATE":
		states.SessionStateReady.Wait()
		var data models.Interaction
		mapstructure.Decode(payload.D, &data)
		app_command.MapInteraction(&data)
	case "GUILD_CREATE":
		
	}
}
