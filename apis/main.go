package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	pb "censys_take_home/grpc"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type grpcServer struct {
	client pb.GRPCClient
}

func main() {
	fmt.Println("Hello1!")

	grpcAddr := os.Getenv("GRPC_SERVER_ADDR") // for docker
	if grpcAddr == "" {
		grpcAddr = "localhost:50000"
	}
	connection, err := grpc.NewClient(grpcAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect to gRPC server:", err)
		return
	}
	defer connection.Close()
	fmt.Println("Connected to gRPC server:", connection)

	client := pb.NewGRPCClient(connection)
	s := &grpcServer{client: client}

	router := gin.Default()
	router.GET("/test", s.printTest)
	router.GET("/getValue/:key", s.getValue)
	router.POST("/setValue", s.setValue)
	router.DELETE("/deleteValue/:key", s.deleteValue)

	router.Run(":8080")
}

func (s *grpcServer) printTest(c *gin.Context) { // defining function like this makes it of type grpcServer, so it can access it
	fmt.Println("In print test")
	resp, err := s.client.GetValue(context.Background(), &pb.GetRequest{Key: "testKey"})
	if err != nil {
		fmt.Println("Error calling GetValue:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call GetValue"})
		return
	}
	fmt.Println("Client: ", resp.Key, resp.Value)
	c.JSON(http.StatusOK, gin.H{"message": "printTest called"})
}

func (s *grpcServer) getValue(c *gin.Context) {
	key := c.Param("key")
	fmt.Println("Get value called with key: " + key)

	resp, err := s.client.GetValue(context.Background(), &pb.GetRequest{Key: key})
	if err != nil {
		fmt.Println("Error calling GetValue:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call GetValue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get called", "key": key, "value": resp.Value})
}

func (s *grpcServer) setValue(c *gin.Context) {
	var body struct {
		Key   string `json:"key" binding:"required"` //required key
		Value string `json:"value" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // wrong input type
		return
	}
	fmt.Println("Set value called with key: " + body.Key + " and value: " + body.Value)

	resp, err := s.client.SetValue(context.Background(), &pb.SetRequest{Key: body.Key, Value: body.Value})
	if err != nil {
		fmt.Println("Error calling SetValue:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call SetValue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Set called", "key": body.Key, "value": body.Value, "response": resp.Result})
}

func (s *grpcServer) deleteValue(c *gin.Context) {
	key := c.Param("key")
	fmt.Println("Delete value called with key: " + key)

	resp, err := s.client.DeleteValue(context.Background(), &pb.DeleteRequest{Key: key})
	if err != nil {
		fmt.Println("Error calling DeleteValue:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call DeleteValue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete called", "key": key, "response": resp.Result})
}
