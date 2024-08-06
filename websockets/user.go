package websockets

import (
	"golang-api/models"
	"golang-api/services"
)

func SendNbUserToAdmin(client *services.Client) error {
	wsData := models.CountStatsUsers(client.Ws)
	return client.Ws.Room("admin").Emit("user:connected", wsData)
}
