package handler

import (
	"fmt"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	services map[string]*grpc.ClientConn
	mutex    sync.RWMutex
}

func NewGRPCClients() *GRPCClients {
	return &GRPCClients{
		services: make(map[string]*grpc.ClientConn),
	}
}

func (g *GRPCClients) AddClient(serviceName, addr string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	log.Printf("order grpc client connected to [%s] grpc server on port [%s]", serviceName, addr)
	g.services[serviceName] = conn
	return nil
}

func (g *GRPCClients) GetClient(serviceName string) (*grpc.ClientConn, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	if conn, ok := g.services[serviceName]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("no connection for service: %s", serviceName)
}

func (g *GRPCClients) CloseAll() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for _, conn := range g.services {
		conn.Close()
	}
}
