package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/smartify/smartify-geofence/pb" // Replace with your actual protobuf package
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	r.GET("/data", func(c *gin.Context) {
		serveWs(c)
	})
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func serveWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}(conn)

	geoFenceServer := NewGeoFenceServer(NewInMemoryGeoFenceRepository())

	go handleIncomingRequests(conn, geoFenceServer)
	go sendFencedLocations(conn, geoFenceServer)

	select {} // Keep the goroutines running
}

func handleIncomingRequests(conn *websocket.Conn, geoFenceServer *GeoFenceServer) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		if messageType == websocket.BinaryMessage {
			addPolygonsRequest := &pb.AddPolygonsRequest{}
			if err := proto.Unmarshal(p, addPolygonsRequest); err != nil {
				log.Println("Failed to parse AddPolygonsRequest:", err)
				return
			}

			for _, polygon := range addPolygonsRequest.Polygons {
				if err := geoFenceServer.Repo.Store(polygon); err != nil {
					log.Println("Error storing polygon:", err)
					return
				}
			}
			log.Printf("Received AddPolygonsRequest with %d polygons.", len(addPolygonsRequest.Polygons))
		}
	}
}

func sendFencedLocations(conn *websocket.Conn, geoFenceServer *GeoFenceServer) {
	for {
		location := getRandomLocation()
		for _, polygon := range geoFenceServer.Repo.GetPolygons() {
			if IsPointInPolygon(location.Latitude, location.Longitude, polygon) {
				fencedLocation := &pb.FencedLocation{
					Latitude:  location.Latitude,
					Longitude: location.Longitude,
					PolygonId: polygon.Id,
				}

				out, err := proto.Marshal(fencedLocation)
				if err != nil {
					log.Println("Failed to encode FencedLocation:", err)
					return
				}

				if err := conn.WriteMessage(websocket.BinaryMessage, out); err != nil {
					log.Println("write:", err)
					return
				}
				log.Printf("Sent FencedLocation for Polygon ID %s.", polygon.Id)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func getRandomLocation() *pb.Location {
	lat := randomNoise(38.9072)
	long := randomNoise(-77.0369)
	return &pb.Location{Latitude: lat, Longitude: long}
}

func randomNoise(coord float64) float64 {
	return coord + rand.Float64()/100
}

func IsPointInPolygon(x, y float64, polygon *pb.Polygon) bool {
	n := len(polygon.Vertices)
	if n < 3 {
		return false
	}

	inside := false
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := polygon.Vertices[i].Latitude, polygon.Vertices[i].Longitude
		xj, yj := polygon.Vertices[j].Latitude, polygon.Vertices[j].Longitude

		intersect := ((yi > y) != (yj > y)) &&
			(x < (xj-xi)*(y-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
	}
	return inside
}
