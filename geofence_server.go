package main

import (
	"github.com/smartify/smartify-geofence/pb"
	"math/rand"
)

type GeoFenceServer struct {
	Repo *InMemoryGeoFenceRepository
}

func NewGeoFenceServer(repo *InMemoryGeoFenceRepository) *GeoFenceServer {
	server := GeoFenceServer{Repo: repo}
	err := repo.Store(
		&pb.Polygon{
			Vertices: []*pb.Point{
				&pb.Point{
					Latitude:  38.9072,
					Longitude: -77.0369,
				},
			},
		})
	if err != nil {
		return nil
	}
	return &server
}

func GetRandomLocation() *pb.Location {
	lat := rand.Float64() * 180
	long := rand.Float64() * 360
	return &pb.Location{Latitude: lat, Longitude: long}
}
