package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "scrapper/grpc"
	"scrapper/service"
)

type server struct {
	pb.UnimplementedScrapperServiceServer
}

func handlePanic() {
	a := recover()

	if a != nil {
		fmt.Printf("Recover from panic: %v", a)
	}
}

func (s *server) FetchVulnerabilities(ctx context.Context, req *pb.VulnerabilityRequest) (*pb.VulnerabilityResponse, error) {
	defer handlePanic()

	query := req.GetName()
	fmt.Printf("New Request: FetchVulnerabilities rpc is called with query value of '%v' \n", query)

	links := service.ExtractVulnerabilitiesLinks(query)

	if len(links) == 0 {
		return nil, fmt.Errorf("there was no matching Vulnerabilities")
	}

	vulnerabilities := service.ExtractVulnerabilitiesDetails(query, links)
	service.SendVulnerabilitiesToDatabase(vulnerabilities)

	result := []*pb.Vulnerability{}
	for _, vulnerability := range vulnerabilities {
		result = append(result, &pb.Vulnerability{
			Name:               vulnerability.Name,
			CVEID:              vulnerability.CVEID,
			PublishedDate:      vulnerability.PublishedDate,
			LastModified:       vulnerability.LastModified,
			Description:        vulnerability.Description,
			VulnerableVersions: vulnerability.VulnerableVersions,
			NVDScore:           vulnerability.NVDScore,
			CNAScore:           vulnerability.CNAScore,
		})
	}

	return &pb.VulnerabilityResponse{
		Vulnerabilities: result,
	}, nil
}

func main() {
	// Initialize RabbitMQ
	service.InitializeRabbitMQ()

	lis, tcpErr := net.Listen("tcp", "0.0.0.0:50051")
	if tcpErr != nil {
		log.Fatalf("Failed to stablish a tcp connections: %v", tcpErr)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterScrapperServiceServer(grpcServer, &server{})

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	go func() {
		fmt.Println("The scrapper server is started ...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	<-channel
	grpcServer.Stop()
	fmt.Println("gRPC server is stoped")
	lis.Close()
	fmt.Println("TCP connection is closed")
}
