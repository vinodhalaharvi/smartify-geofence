package main

import (
	"github.com/smartify/smartify-geofence/pb"
	"math/rand"
)

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
