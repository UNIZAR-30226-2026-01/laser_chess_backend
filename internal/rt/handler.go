package rt

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// fichero que se encarga de upgradear conexiones HTTP a websoket
// tendrá endpoints para partidas competitivas y privadas.

func UpgradeConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error){
	var upgrader *websocket.Upgrader
	conn, err := upgrader.Upgrade(w, r, nil)
	return conn, err
}

//TODO: todo
