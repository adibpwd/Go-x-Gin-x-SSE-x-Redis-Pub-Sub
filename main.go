//main.go
package main

import (
	"context"
	"fmt"
	"io"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: "172.30.0.5:6379",
})

func main() {
	app := gin.Default()

	app.GET("/pub-test/:msg", func(c *gin.Context) {
		err := rdb.Publish(ctx, "mychannel1", c.Param("msg")).Err()
		if err != nil {
			fmt.Println(err)
		}
	})
	app.GET("/sub-test-sse/:name", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		name := c.Param("name")
		pubsub := rdb.Subscribe(ctx, "mychannel1")

		defer pubsub.Close()

		c.Stream(func(w io.Writer) bool {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Println(err)
				return false
			}
			fmt.Println("mantaps")
			fmt.Println(msg.Channel, msg.Payload)

			c.SSEvent("message", msg.Payload+" ><><><> "+name)
			return true
		})
	})
	app.Run(":3000")
}
