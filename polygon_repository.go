package main

import (
	"github.com/smartify/smartify-geofence/pb"
	"sync"
)

type GeoFenceRepository interface {
	Store(polygon *pb.Polygon) error
	GetPolygons() []*pb.Polygon
}

type InMemoryGeoFenceRepository struct {
	Polygons []*pb.Polygon
	mu       sync.Mutex
}

func NewInMemoryGeoFenceRepository() *InMemoryGeoFenceRepository {
	return &InMemoryGeoFenceRepository{Polygons: []*pb.Polygon{}, mu: sync.Mutex{}}
}

func (r *InMemoryGeoFenceRepository) Store(p *pb.Polygon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Polygons = append(r.Polygons, p)
	return nil
}

func (r *InMemoryGeoFenceRepository) GetPolygons() []*pb.Polygon {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Polygons
}
