package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	locator "github.com/hunkvm/locator/client"
	hello "github.com/hunkvm/locator/example/api"
	hellosrv "github.com/hunkvm/locator/example/server"

	"github.com/hunkvm/locator/pkg/types"
	"github.com/hunkvm/locator/pkg/validation"

	health "google.golang.org/grpc/health/grpc_health_v1"

	healthsrv "google.golang.org/grpc/health"

	"google.golang.org/grpc"
)

func main() {
	address := flag.String("address", ":50051", "gRPC server listen address")
	registry := flag.String("registry", "localhost:8881", "Registry address")
	name := flag.String("name", "hello", "Service name to register")
	flag.Parse()

	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	healthServer := healthsrv.NewServer()
	helloServer := hellosrv.NewServer()

	health.RegisterHealthServer(grpcServer, healthServer)
	hello.RegisterHelloServiceServer(grpcServer, helloServer)

	healthServer.SetServingStatus(*name, health.HealthCheckResponse_SERVING)

	log.Printf("gRPC server listening on %s\n", *address)
	go registerService(*address, *registry, *name)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func registerService(serviceAddr string, registryAddr string, serviceName string) {
	if e := validation.ValidateAddress(registryAddr); e != nil {
		panic(e)
	}
	if e := validation.ValidateAddress(serviceAddr); e != nil {
		panic(e)
	}

	client, err := locator.NewLocatorClient(
		locator.WithAddress(registryAddr),
		locator.WithTimeout(5*time.Second),
	)
	if err != nil {
		panic(err)
	}

	service := types.Service{
		Name:    serviceName,
		Address: serviceAddr,
		Enabled: true,
		Version: "v1",
		Metadata: map[string]string{
			"Namespace": "default",
			"Region":    "southeast-asia",
			"Zone":      "indochina",
			"SubZone":   "hochiminh",
		},
		HealthCheck: &types.HealthConfig{
			Service: serviceName,
			Address: serviceAddr,
		},
	}
	id, err := client.Registry.Register(context.Background(), service)
	if err != nil {
		panic(err)
	}
	fmt.Println("Service registered with ID:", id)
	services, err := client.Registry.Retrieve(
		context.Background(),
		map[string][]any{"ID": {id}},
	)
	if err != nil {
		panic(err)
	}
	for _, service := range services {
		fmt.Printf("Retrieved service: %+v\n", service)
	}
}
