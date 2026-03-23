package rt

// fichero que se encarga de upgradear conexiones HTTP a websoket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // size of read buffer
	WriteBufferSize: 1024, // size of write buffer
	CheckOrigin: func(r *http.Request) bool { //allowing CORS request

		//TODO: cambiar esto a algo seguro
		return true
	},
}

func UpgradeConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}
