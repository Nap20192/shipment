package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Nap20192/shipment/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewShipmentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Println("--- Creating Shipment ---")
	createResp, err := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
		Origin:      "Warsaw, PL",
		Destination: "Berlin, DE",
		Details: &pb.ShipmentDetails{
			Weight:          150.5,
			DimensionLength: 100,
			DimensionWidth:  80,
			DimensionHeight: 60,
		},
		DriverDetails: &pb.DriverDetails{
			Name: "John Doe",
		},
	})
	if err != nil {
		log.Fatalf("could not create shipment: %v", err)
	}
	shipmentID := createResp.GetShipment().GetId()
	fmt.Printf("Created Shipment ID: %s\n", shipmentID)

	fmt.Println("\n--- Updating Status to IN_TRANSIT ---")
	_, err = client.UpdateShipmentStatus(ctx, &pb.UpdateShipmentStatusRequest{
		Id:        shipmentID,
		NewStatus: pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT,
	})
	if err != nil {
		log.Fatalf("could not update status: %v", err)
	}
	fmt.Println("Status updated to IN_TRANSIT")

	fmt.Println("\n--- Updating Status to DELIVERED ---")
	_, err = client.UpdateShipmentStatus(ctx, &pb.UpdateShipmentStatusRequest{
		Id:        shipmentID,
		NewStatus: pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED,
	})
	if err != nil {
		log.Fatalf("could not update status: %v", err)
	}
	fmt.Println("Status updated to DELIVERED")

	fmt.Println("\n--- Fetching Shipment History ---")
	historyResp, err := client.GetShipmentEventHistory(ctx, &pb.GetShipmentEventHistoryRequest{
		ShipmentId: shipmentID,
	})
	if err != nil {
		log.Fatalf("could not get history: %v", err)
	}

	fmt.Printf("History for shipment %s:\n", shipmentID)

	for _, event := range historyResp.GetEvents() {
		fmt.Printf("- Event: %s, Time: %s, Description: %s\n",
			event.GetEventName(),
			event.GetCreatedAt().AsTime().Format(time.RFC3339),
			string(event.GetPayload()))
	}
}
