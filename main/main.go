package main

import (
	"context"
	"fmt"
	"log"

	pb "goscrapper/main/grpc"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	app := fiber.New()

	clientConn, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to the gRPC server: %v \n", err)
	}
	defer clientConn.Close()
	grpcClient := pb.NewScrapperServiceClient(clientConn)

	app.Get("/api/vulnerabilities/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		fmt.Printf("New Request: /vulnerabilities/:%v \n", name)
		res, err := grpcClient.FetchVulnerabilities(context.Background(), &pb.VulnerabilityRequest{
			Name: name,
		})
		if err != nil {
			fmt.Printf("Failed to call FetchVulnerabilities rpc: %v \n", err)
			return err
		}

		return c.JSON(res.Vulnerabilities)
	})

	app.Listen(":3040")
}
