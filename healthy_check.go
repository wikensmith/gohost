package gohost

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"net/http"
	"time"
)

var errChan chan *amqp.Error
var isConnectClosed = false // 连接中断 ， 默认是false， 即没有中断
var notifyClose chan *amqp.Error

// 健康检测
func healthyCheck(c *gin.Context) {
	if isConnectClosed {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "连接已关闭",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "连接正常",
		})
	}
}

func checkClose() {
	//for {
	//	if v, open := <-notifyClose; !open {
	//		fmt.Println("连接已断开", v.Reason)
	//		isConnectClosed = true
	//		Start()
	//		break
	//	} else {
	//		fmt.Println("bbbbbbbb:", v.Reason)
	//	}
	//	fmt.Println("aaaaaaa")
	//}
	for {
		time.Sleep(time.Second * 5)
		if conn == nil {
			isConnectClosed = true
			fmt.Println("连接断开，重新连接")
			conn, _ = GetConnection()
			if conn != nil {
				if !conn.IsClosed() {
					_ = conn.Close()
					Start()
				}
			}
			continue
		}
		if conn.IsClosed() {
			isConnectClosed = true
			fmt.Println("连接断开，重新连接")
			conn, _ = GetConnection()
			if conn != nil {
				if !conn.IsClosed() {
					_ = conn.Close()
					Start()
				}
			}
			continue
		} else {
			isConnectClosed = false
		}
	}
}

func HealthyCheck() {
	r := gin.Default()
	r.GET("/healthyCheck", healthyCheck)
	_ = r.Run("0.0.0.0:" + Params.HealthyPort)
}
