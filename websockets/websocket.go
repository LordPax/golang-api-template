package websockets

import (
	"context"
	"fmt"
	"golang-api/middlewares"
	"golang-api/models"
	"net/http"

	"github.com/LordPax/sockevent"
	"github.com/gin-gonic/gin"
)

func RegisterWebsocket(r *gin.Engine) {
	ws := sockevent.GetWebsocket()

	ws.OnConnect(connect)
	ws.OnDisconnect(disconnect)
	ws.On("ping", func(client *sockevent.Client, message any) error {
		logged := client.Get("logged").(bool)
		if !logged {
			fmt.Printf("Client %s sent message: %v\n", client.ID, message)
			return client.Emit("pong", "pong")
		}

		user := client.Get("user").(models.User)
		fmt.Printf("Client %s (%s) sent message: %v\n", client.ID, user.Username, message)
		return client.Emit("pong", "pong")
	})

	r.GET("/ws",
		middlewares.IsLoggedIn(false),
		func(c *gin.Context) {
			connectedUser, ok := c.Get("connectedUser")
			if !ok {
				ws.WsHandler(c.Writer, c.Request)
				return
			}

			ctx := context.WithValue(c.Request.Context(), "connectedUser", connectedUser)
			ws.WsHandler(c.Writer, c.Request.WithContext(ctx))
		},
	)
}

func connect(client *sockevent.Client, wr http.ResponseWriter, r *http.Request) error {
	connectedUser := r.Context().Value("connectedUser")
	ok := connectedUser != nil && connectedUser.(models.User).ID != 0
	client.Set("logged", ok)

	if err := SendNbUserToAdmin(client); err != nil {
		return err
	}

	if !ok {
		fmt.Printf("Client %s connected, len %d\n",
			client.ID,
			len(client.Ws.GetClients()),
		)
		return client.Emit("connected", nil)
	}

	user := connectedUser.(models.User)
	client.Set("user", user)

	if user.IsRole(models.ROLE_ADMIN) {
		client.Ws.Room("admin").AddClient(client)
	}

	fmt.Printf("Client %s (%s) connected, len %d\n",
		client.ID,
		user.Username,
		len(client.Ws.GetClients()),
	)

	return client.Emit("connected", nil)
}

func disconnect(client *sockevent.Client) error {
	logged := client.Get("logged").(bool)

	if err := SendNbUserToAdmin(client); err != nil {
		return err
	}

	if !logged {
		fmt.Printf("Client %s disconnected, len %d\n",
			client.ID,
			len(client.Ws.GetClients()),
		)
		return nil
	}

	user := client.Get("user").(models.User)

	fmt.Printf("Client %s (%s) disconnected, len %d\n",
		client.ID,
		user.Username,
		len(client.Ws.GetClients()),
	)

	return nil
}
