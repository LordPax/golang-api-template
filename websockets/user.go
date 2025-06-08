package websockets

import (
	"golang-api/models"

	"github.com/LordPax/sockevent"
)

func SendNbUserToAdmin(client *sockevent.Client) error {
	wsData := models.CountStatsUsers(client.Ws)
	return client.Ws.Room("admin").Emit("user:connected", wsData)
}
