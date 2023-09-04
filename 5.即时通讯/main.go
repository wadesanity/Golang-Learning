package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"talk/c1"
	"talk/util"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	util.Logger.Debugln("Start main")
	c1.InitHub()
	engine := gin.New()
	engine.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			util.Logger.Errorf("upgrade err:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "升级ws失败",
			})
			return
		}
		idString, ok := c.GetQuery("id")
		if !ok {
			util.Logger.Errorf("query id miss")
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  "升级ws失败",
			})
			return
		}
		id, err := strconv.Atoi(idString)
		if err != nil {
			util.Logger.Errorf("query id .cannot convert to int")
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  "升级ws失败",
			})
			return
		}

		client := c1.NewClientInterface(conn, id)
		client.Run()
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "结束",
		})
	})
	engine.Run(":8080")
}
