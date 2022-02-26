package gateway

import (
	"bothoi/models"

	"github.com/mitchellh/mapstructure"
)

func dispatchHandler(payload models.GatewayPayload) {
	switch payload.T {
	case "READY":
		mapstructure.Decode(payload.D, &ActiveSessionState)
		SessionReady.Done()
	case "GUILD_CREATE":
	case "INTERACTION_CREATE":
	}
}
