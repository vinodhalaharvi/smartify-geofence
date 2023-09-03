package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/smartify/smartify-geofence/pb"
	"google.golang.org/protobuf/proto"
	"log"
	"math/rand"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GeoFenceServer struct {
	Repo GeoFenceRepository
}

func NewGeoFenceServer(
	repo GeoFenceRepository,
) *GeoFenceServer {
	server := GeoFenceServer{Repo: repo}
	server.Repo = repo
	err := repo.Store(&pb.Polygon{
		Id: "1",
		Vertices: []*pb.Point{
			{Latitude: 38.9072, Longitude: -77.0369},
		},
	})
	if err != nil {
		return nil
	}
	return &server
}

func main() {
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		serveData(c)
	})
	r.Run(":8080")
}

func serveData(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	geoFenceServer := NewGeoFenceServer(NewInMemoryGeoFenceRepository())

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error while reading message:", err)
			return
		}

		if messageType == websocket.BinaryMessage {
			// Handle AddPolygonsRequest
			addPolygonsRequest := &pb.AddPolygonsRequest{}
			if err := proto.Unmarshal(p, addPolygonsRequest); err != nil {
				log.Println("Failed to parse AddPolygonsRequest:", err)
				return
			}

			// Store or process the polygons
			for _, polygon := range addPolygonsRequest.Polygons {
				if err := geoFenceServer.Repo.Store(polygon); err != nil {
					log.Println("Error storing polygon:", err)
					return
				}
			}

			log.Printf("Received polygons")
		}

		// Continue sending FencedLocation or whatever you are sending
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
					log.Println("Write error:", err)
					return
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// getRandomLocation function for demonstration
func getRandomLocation() *pb.Location {
	lat := randomNoise(38.9072)
	long := randomNoise(-77.0369)
	return &pb.Location{Latitude: lat, Longitude: long}
}

func randomNoise(coord float64) float64 {
	return coord + rand.Float64()/100
}

// IsPointInPolygon checks if a point (x, y) is inside a given polygon
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
