package main

import (
	"context"
	"fmt"
	"greet/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Hello I am Client")

	// With SSL
	tls := true
	opts := grpc.WithInsecure()

	if tls {
		certFile := "greet/ssl/ca.crt" // Certificate Authority
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificates : %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not conect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	// fmt.Printf("Created Client %v", c)

	doUnary(c)
	// CalculateSum(c)
	// doServerStreaming(c)
	// doServerStreamingPrimeDecomposition(c)
	// doClientStreaming(c)
	// doClientStreamingComputeAverage(c)
	// doBidiStreaming(c)
	// doBidiFindMaxStreaming(c)
	// doErrorUnary(c)
	// doUnaryWithDeadline(c, 5*time.Second) // Will Pass
	// doUnaryWithDeadline(c, 1*time.Second) // Will Timeout

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a Unary RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Deepak",
			LastName:  "Mahana ",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC : %v", err)
	}

	log.Printf("Response from Greet : %v", res.Result)
}

// Exercise
func CalculateSum(c greetpb.GreetServiceClient) {

	fmt.Println("Calculate Sum using Unary RPC")

	req := &greetpb.CalculateRequest{
		Calvalue: &greetpb.CalculateValues{
			X: 10,
			Y: 10,
		},
	}

	res, err := c.CalculateSum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC : %v", err)
	}

	log.Printf("Response from Greet : %v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Deepak",
			LastName:  "Mahana",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While Calling GreetManyTimes RPC : %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes : %v", msg.GetResult())
	}

}

func doServerStreamingPrimeDecomposition(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While Calling PrimeNumberDecomposition RPC : %v", err)
	}

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			// reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Println(res.GetPrimeFactor())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cuba",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet : %v", err)
	}

	// Iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending Req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error While Recieving Response from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)

}

func doClientStreamingComputeAverage(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a client streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream : %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending Number : %v\n", number)
		stream.Send(&greetpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while Receiving Response : %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())

}

func doBidiStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting BiDi Streaming Client RPC")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream : %v", err)
		return
	}

	// Create Data to send
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deepak",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cuba",
			},
		},
	}

	// Create a channel
	waitc := make((chan struct{}))

	// Send bunch of messages to the client (GoRoutine)
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// We recieve a bunch of messages from the client (GoRoutine)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving : %v", err)
				break
			}
			fmt.Printf("Received : %v\n", res.GetResult())
		}
	}()

	// Block until everything is done
	<-waitc

}

func doBidiFindMaxStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a FindMaximum Bidi Streaming RPC....")
	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum : %v", err)
	}

	waitc := make(chan struct{})

	// Send go routine
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			stream.Send(&greetpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// Receive go routing
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received a new maximum of %v\n", maximum)
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do a SquareRoot Unary RPC")

	// Correct Call
	doErrorCall(c, 10)

	// Error Call
	doErrorCall(c, -2)
}

func doErrorCall(c greetpb.GreetServiceClient, n int32) {

	res, err := c.SquareRoot(context.Background(), &greetpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, OK := status.FromError(err)
		if OK {
			// Actual error from gRPC (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number !")
			}
		} else {
			log.Fatalf("Big Error Calling SquareRoot : %v", err)
		}
	}
	fmt.Printf("Result of Square Root of %v : %v\n", n, res.GetNumberRoot())
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {

	fmt.Println("Starting to do a UnaryWithDeadline RPC...")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Deepak",
			LastName:  "Mahana ",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, OK := status.FromError(err)
		if OK {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit ! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected Error : %v", statusErr)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline RPC : %v", err)
		}
		return
	}
	log.Printf("Response from Greet : %v", res.Result)
}
