package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/smartify/smartify-geofence/pb" // make sure this import path is correct
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
			log.Printf("Error closing connection: %v", err)
		}
	}(conn)

	for {
		// Simulate generating FencedLocation data (replace this with real data)
		fencedLocation := &pb.FencedLocation{
			Latitude:  38.9072 + rand.Float64(),
			Longitude: -77.0369 + rand.Float64(),
			PolygonId: "0",
		}

		out, err := proto.Marshal(fencedLocation)
		if err != nil {
			log.Println("Failed to encode data:", err)
			return
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, out); err != nil {
			log.Println("Write error:", err)
			return
		}

		// Add a delay or a trigger for new data
		time.Sleep(10 * time.Second)
	}
}

func main() {
	r := gin.Default()

	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	r.Use(cors.New(config))

	r.Static("/static", "./static")
	r.GET("/data", serveData)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
