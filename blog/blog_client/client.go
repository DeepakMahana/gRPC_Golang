package main

import (
	"context"
	"fmt"
	"greet/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Blog Client")
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	// 1. Create Blog
	fmt.Println("Create a Blog")
	blog := &blogpb.Blog{
		AuthorId: "Deepak",
		Title:    "My First Blog",
		Content:  "Content of the first Blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected Error : %v", err)
	}
	fmt.Printf("Blog has been created %v", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// 2. Read Blog
	fmt.Println("Reading Blog")
	_, readErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "60eaf2ca7626531bf703318a"})
	if readErr != nil {
		fmt.Printf("Error happened while reading : %v \n", readErr)
	}
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading : %v \n", err)
	}
	fmt.Printf("Blog was read : %v \n", readBlogRes)

	// 3. Update Blog
	fmt.Println("Update a Blog")
	updateBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Deepak Mahana",
		Title:    "My Updated Blog",
		Content:  "Content of the updated Blog",
	}
	updateBlogRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: updateBlog})
	if updateErr != nil {
		log.Fatalf("Unexpected Error : %v \n", updateErr)
	}
	fmt.Printf("Blog has been updated %v \n", updateBlogRes)

	// 4. Delete Blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Printf("Error happened while deleting : %v \n", updateErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	// 5. List Blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog RPC : %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}
