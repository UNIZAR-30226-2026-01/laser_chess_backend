package test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var access_token string = ""
var session_cookies []*http.Cookie
var router *gin.Engine

type LoginResponse struct {
	Token string `json:"access_token"`
}

func init() {
	godotenv.Load("../.env")

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatalln("Error: DB_URL no encontrada")
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalln("Error: No se pudo conectar con la base de datos")
	}
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalln("Error: No hay conexión real con la DB:", err)
	}

	store := db.NewStore(dbPool)

	matchRegistry := rt.NewMatchRegistry()

	privateHub := rt.NewPrivateHub()

	publicHub := rt.NewPublicHub()

	router = api.SetupRouter(store, matchRegistry, privateHub, publicHub)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func performRequest(method, path string, body []byte, token string, ctx context.Context) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	for _, cookie := range session_cookies {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := http.Response{Header: w.Header()}
	newCookies := resp.Cookies()
	if len(newCookies) > 0 {
		session_cookies = newCookies
	}

	return w
}

func LoginUser(t *testing.T, username string) *httptest.ResponseRecorder {
	session_cookies = nil

	password := username + username
	datos := map[string]string{
		"credential": username,
		"password":   password,
	}
	body, _ := json.Marshal(datos)

	w := performRequest("POST", "/login", body, "", nil)

	if w.Code >= 200 && w.Code < 300 {
		var lr LoginResponse
		json.Unmarshal(w.Body.Bytes(), &lr)
		access_token = lr.Token
	}
	return w
}

func TestLogin(t *testing.T) {
	w := LoginUser(t, "user1")
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Se esperaba 2xx, pero recibimos %d", w.Code)
	}
}

func TestRegister(t *testing.T) {
	body := []byte(`{"username": "user151", "password": "user151user151", "mail": "user151@gmail.com"}`)
	w := performRequest("POST", "/register", body, "", nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestRefresh(t *testing.T) {
	TestLogin(t)
	body := []byte(`{"username": "user151", "password": "user151user151", "mail": "user151@gmail.com"}`)
	w := performRequest("POST", "/refresh", body, "", nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestLogout(t *testing.T) {
	TestLogin(t)
	w := performRequest("POST", "/logout", []byte{}, "", nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetOwnAccount(t *testing.T) {
	w := performRequest("GET", "/api/account/", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetOtherAccount(t *testing.T) {
	w := performRequest("GET", "/api/account/1", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetXPInfo(t *testing.T) {
	w := performRequest("GET", "/api/account/xp", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestUpdateAccount(t *testing.T) {
	w := performRequest("POST", "/api/account/update", []byte(`{}`), access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestDeleteAccount(t *testing.T) {
	LoginUser(t, "user151")
	w := performRequest("DELETE", "/api/account/delete", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
	LoginUser(t, "user1")
}

func TestGetMatch(t *testing.T) {
	w := performRequest("GET", "/api/match/2", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetUserHistory(t *testing.T) {
	w := performRequest("GET", "/api/match/history/2", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestCreateItemOwner(t *testing.T) {
	body := []byte(`{"item_id": 11}`)
	w := performRequest("POST", "/api/item", body, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetUserItems(t *testing.T) {
	w := performRequest("GET", "/api/item/inventory", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetAllElos(t *testing.T) {
	w := performRequest("GET", "/api/rating/1", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetTopRankUsers(t *testing.T) {
	w := performRequest("GET", "/api/rating/top/blitz", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestCreateFriendship(t *testing.T) {
	body := []byte(`{"receiver_username": "user150"}`)
	w := performRequest("POST", "/api/friendship", body, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetUserFriendships(t *testing.T) {
	w := performRequest("GET", "/api/friendship", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestAcceptFriendship(t *testing.T) {
	LoginUser(t, "user150")
	w := performRequest("PUT", "/api/friendship/user1", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
	LoginUser(t, "user1")
}

func TestDeleteFriendship(t *testing.T) {
	w := performRequest("DELETE", "/api/friendship/user3", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestEventHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	w := performRequest("GET", "/api/events", nil, access_token, ctx)

	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestRegisterDevice(t *testing.T) {
	body := []byte(`{"token": "aaaa"}`)
	w := performRequest("POST", "/api/device/register", body, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestGetWSTicket(t *testing.T) {
	w := performRequest("GET", "/api/rt/ticket", nil, access_token, nil)
	if w.Code >= 300 || w.Code < 200 {
		t.Errorf("Expected 2xx, received %d", w.Code)
	}
}

func TestRTChallenge(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()
	wsURL := strings.Replace(ts.URL, "http", "ws", 1)

	LoginUser(t, "user1")
	token1 := access_token

	LoginUser(t, "user2")
	token2 := access_token

	// Hacer friend request de user1 a user2
	body := []byte(`{"receiver_username": "user2"}`)
	w := performRequest("POST", "/api/friendship", body, token1, nil)
	if w.Code != 201 {
		t.Fatalf("Error creando friendship: %d", w.Code)
	}

	// Aceptar con user2
	w = performRequest("PUT", "/api/friendship/user1", nil, token2, nil)
	if w.Code != 200 {
		t.Fatalf("Error aceptando friendship: %d", w.Code)
	}

	headers1 := http.Header{}
	headers1.Add("Authorization", "Bearer "+token1)

	conn1, resp1, err := websocket.DefaultDialer.Dial(wsURL+"/api/rt/challenge?username=user2&board=0&starting_time=300&time_increment=5", headers1)
	if err != nil {
		status := 0
		if resp1 != nil {
			status = resp1.StatusCode
		}
		t.Fatalf("Error en conexion WS de user1: %v, status HTTP: %d", err, status)
	}
	defer conn1.Close()

	headers2 := http.Header{}
	headers2.Add("Authorization", "Bearer "+token2)

	conn2, resp2, err := websocket.DefaultDialer.Dial(wsURL+"/api/rt/challenge/accept?username=user1", headers2)
	if err != nil {
		status := 0
		if resp2 != nil {
			status = resp2.StatusCode
		}
		t.Fatalf("Error en conexion WS de user2: %v, status HTTP: %d", err, status)
	}
	defer conn2.Close()
}

func TestRTGoIntoMatchmaking(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()
	wsURL := strings.Replace(ts.URL, "http", "ws", 1)

	LoginUser(t, "user3")
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+access_token)

	conn, resp, err := websocket.DefaultDialer.Dial(wsURL+"/api/rt/matchmaking?board=0&time_base=300&time_increment=5", headers)
	if err != nil {
		status := 0
		if resp != nil {
			status = resp.StatusCode
		}
		t.Fatalf("Error en conexion WS de user3: %v, status HTTP: %d", err, status)
	}
	conn.Close()

	conn, resp, err = websocket.DefaultDialer.Dial(wsURL+"/api/rt/matchmaking?board=0&time_base=300&time_increment=5", headers)
	if err != nil {
		status := 0
		if resp != nil {
			status = resp.StatusCode
		}
		t.Fatalf("Error en segunda conexion WS: %v, status HTTP: %d", err, status)
	}
	conn.Close()
}
