package server

import (
	"BaseWebSocket/domain/authenticate/delivery"
	_authenticationRepo "BaseWebSocket/domain/authenticate/repository"
	_authenticationUseCase "BaseWebSocket/domain/authenticate/usecase"
	"BaseWebSocket/service/server/client"
	"BaseWebSocket/service/server/lib"
	"BaseWebSocket/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"github.com/spf13/viper"
)

var usableClientID uint32
var usableClientIDLock sync.Mutex

var authenticationDelivery *lib.ApiRouter
var connectionTimeOut time.Duration

var serverConfig *viper.Viper

var globalHub *client.Hub

// 取得client ID
func getUsableClientID() uint32 {
	usableClientIDLock.Lock()
	defer usableClientIDLock.Unlock()

	usableClientID++
	if usableClientID == 0 {
		usableClientID++
	}
	return usableClientID
}

type ResponseError struct {
	Message string `json:"message"`
}

type DeliverServer struct {
}

func initRun() {
	serverConfig = utils.SetConfigPath()
}

// 初始化 server 使用的repo & useCase & router
func NewDeliverServer() *DeliverServer {
	initRun()
	// mysqlSettingWriteConn := utils.NewMysql(serverConfig, "setting", "master")
	// mysqlSettingReadConn := utils.NewMysql(serverConfig, "setting", "slave")
	settingRedisConn := utils.NewRedis(serverConfig, "setting")

	// restyClient := resty.New()

	connectionTimeOut = serverConfig.GetDuration("connectionTimeOut")

	authenticRedisRepo := _authenticationRepo.NewRedisRepository(settingRedisConn)
	authenticUseCase := _authenticationUseCase.NewAuthUseCase(authenticRedisRepo)

	authenticationDelivery = delivery.NewRouter(authenticUseCase)
	handler := &DeliverServer{}

	globalHub = client.NewHub()
	go globalHub.Run()

	lib.ApiControllerInit(authenticationDelivery)

	return handler
}

// websocket server 啟動
func (ds *DeliverServer) Start() {
	srv := &http.Server{
		Addr: ":8999",
	}
	http.HandleFunc("/webSocketServer", WsHandler)
	fmt.Println("Server start at port 8999")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Listen And Serve Fail")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown Server Fail")
	}

	globalHub.ServerMaintain([]byte{})
	time.Sleep(60 * time.Second)
	fmt.Println("Exiting")
}

// websocketHandler
func WsHandler(w http.ResponseWriter, r *http.Request) {
	upGrader := &websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: time.Second * 5,
	}

	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade fail")
		return
	}

	newClientID := getUsableClientID()
	newClient := lib.NewClientObj(globalHub, newClientID, conn, connectionTimeOut)
	err = newClient.Conn.SetReadDeadline(time.Now().Add(connectionTimeOut * time.Second))
	if err != nil {
		fmt.Println("SetReadDeadline fail")
	}

	go lib.ReadPeerMessage(newClient)
	go lib.WritePeerMessage(newClient)

	newClient.Hub.Register <- newClient
}
