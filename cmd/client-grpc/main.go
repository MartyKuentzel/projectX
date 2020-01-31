package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	v1 "github.com/MartyKuentzel/projectX/pkg/api/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewProductServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	date, _ := ptypes.TimestampProto(t)

	// Call Create
	req1 := v1.CreateRequest{
		Api: apiVersion,
		Product: &v1.ProductProto{
			Name:        "Potato",
			Price:       "5€",
			Creator:     "Marty",
			Unit:        "Kg",
			Description: "Buy my Potato",
			Category:    "vegetable",
			Date:        date,
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	id := res1.Id

	// Read
	req2 := v1.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	log.Printf("Read result: <%+v>\n\n", res2)

	// Update
	req3 := v1.UpdateRequest{
		Api: apiVersion,
		Product: &v1.ProductProto{
			Id:          res2.Product.Id,
			Name:        res2.Product.Name,
			Price:       res2.Product.Price,
			Creator:     res2.Product.Creator + " + updated",
			Unit:        res2.Product.Unit,
			Description: res2.Product.Description + " + updated",
			Category:    res2.Product.Category,
			Date:        res2.Product.Date,
		},
	}

	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Update result: <%+v>\n\n", res3)

	// Call ReadAll
	req4 := v1.ReadAllRequest{
		Api: apiVersion,
	}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatalf("ReadAll failed: %v", err)
	}
	log.Printf("ReadAll result: <%+v>\n\n", res4)

	// // Delete
	// req5 := v1.DeleteRequest{
	// 	Api: apiVersion,
	// 	Id:  id,
	// }
	// res5, err := c.Delete(ctx, &req5)
	// if err != nil {
	// 	log.Fatalf("Delete failed: %v", err)
	// }
	// log.Printf("Delete result: <%+v>\n\n", res5)
}
