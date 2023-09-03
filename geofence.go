package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/smartify/smartify-geofence/pb" // Replace with your actual protobuf package
	"google.golang.org/protobuf/proto"
	"log"
	"math/rand"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	repo := NewInMemoryGeoFenceRepository()
	server := NewGeoFenceServer(repo)

	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		HandleConnections(c.Writer, c.Request, server)
	})
	err := r.Run(":8080")
	if err != nil {
		log.Fatalln("Could not start server:", err)
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request, server *GeoFenceServer) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Could not close connection: %v", err)
		}
	}(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if messageType == websocket.BinaryMessage {
			addPolygonsReq := &pb.AddPolygonsRequest{}
			if err := proto.Unmarshal(p, addPolygonsReq); err != nil {
				log.Println("Could not unmarshal AddPolygonsRequest:", err)
			} else {
				log.Println("Received AddPolygonsRequest")
				for _, polygon := range addPolygonsReq.Polygons {
					err := server.Repo.Store(polygon)
					if err != nil {
						log.Printf("Could not store polygon: %v", err)
					}
				}
			}
		}

		// Send fenced location data
		fencedLocation := &pb.FencedLocation{ /* populate fields here */ }
		data, err := proto.Marshal(fencedLocation)
		if err != nil {
			log.Println("Could not marshal FencedLocation:", err)
		} else {
			log.Println("Sending FencedLocation")
			if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Println("Could not send FencedLocation:", err)
			}
		}
	}
}

type GeoFenceServer struct {
	Repo *InMemoryGeoFenceRepository
}

func NewGeoFenceServer(repo *InMemoryGeoFenceRepository) *GeoFenceServer {
	return &GeoFenceServer{Repo: repo}
}

func GetRandomLocation() *pb.Location {
	lat := rand.Float64() * 180
	long := rand.Float64() * 360
	return &pb.Location{Latitude: lat, Longitude: long}
}
