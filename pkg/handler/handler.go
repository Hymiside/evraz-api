package handler

import (
	"context"
	"fmt"
	"github.com/Hymiside/evraz-api/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

// Получаю канал

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s]	REQUEST: %s %s    STATUS-CODE: %d    LATENSY: %s\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	}))

	router.GET("/trends")
	router.GET("/ws", h.wsHandler)

	return router
}

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) wsHandler(c *gin.Context) {
	fmt.Println("hellosdf")

	conn, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := h.services.Register("hello")
	fmt.Printf("chan: %v\n", ch)
	defer h.services.Unregister("hello")

	conn.SetCloseHandler(func(code int, text string) error {
		cancel()
		return nil
	})

	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			return
	//		default:
	//			_, p, err := conn.ReadMessage()
	//			if err != nil {
	//				log.Println(err)
	//				return
	//			}
	//			conn.WriteJSON(string(p))
	//		}
	//	}
	//}()

	for {
		select {
		case d := <-ch.Ch:
			fmt.Println(d)

		case <-ctx.Done():
			return
		}
	}
}
