package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/smartify/smartify-livetable/pb"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

func serveData(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	fiveNames := []string{"Alice", "Bob", "Carol", "Dave", "Eve"}
	ages := []int32{30, 40, 50, 60, 70}

	for {
		// Generate table data (replace with real data as needed)
		// cycle through the names
		for i := 0; i < len(fiveNames); i++ {
			rows := []*pb.Row{
				{
					Id:   int64(rand.Int()),
					Name: fiveNames[i],
					Age:  ages[i],
				},
			}
			data := &pb.TableData{Rows: rows}

			out, err := proto.Marshal(data)
			if err != nil {
				log.Println("Failed to encode data:", err)
				return
			}

			if err := conn.WriteMessage(websocket.BinaryMessage, out); err != nil {
				log.Println("Write error:", err)
				return
			}

			// Add a delay or a trigger for new data
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	r := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Allow all origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	r.Use(cors.New(config))

	r.Static("/static", "./static")
	r.GET("/data", serveData)

	// Start the server
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
