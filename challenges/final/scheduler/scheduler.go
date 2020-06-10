package scheduler

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/CodersSquad/dc-labs/challenges/final/controller"

	pb "github.com/CodersSquad/dc-labs/challenges/final/proto"
	"google.golang.org/grpc"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

type Job struct {
	Address string
	RPCName string
	Info [4]string
}

var counter int


func schedule(job Job, name string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTaskClient(conn)
	controller.ChangeStatus(name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch job.RPCName {
	case "test":
		r, err := c.SayHello(ctx, &pb.TestRequest{Name: job.RPCName})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Scheduler: RPC respose from %s : %s", job.Address, r.GetMessage())
		controller.Register(name, counter)
	case "image":
		controller.RegisterImage(name, job.Info[0], job.Info[2], counter)
		img := pb.Image{
			Workload: job.Info[2], 
			Index: int64(controller.Results[job.Info[2]].LastIndex), 
			Filepath: job.Info[0],
			Filter: job.Info[3],
		}
		r, err := c.FilterImg(ctx, &pb.ImgRequest{Name: "Image Filter", Img: &img })
		if err != nil {
			log.Fatalf("could not proccess image: %v", err)
		}
		log.Printf("Scheduler: RPC respose from %s : %s was filtered", job.Address, r.GetMessage())
		
	}
	controller.ChangeStatus(name)
	counter++
}

func Start(jobs chan Job) error {
	counter = 0
	for {
		job := <-jobs
		time.Sleep(time.Second * 5)
		lowestUsage := 99999
		lowestPort := 0
		worker := controller.Worker{}
		for _, data := range controller.Nodes {
			if data.Usage < lowestUsage {
				lowestPort = data.Port
				lowestUsage = data.Usage
				worker = data
			}
		}
		controller.IncreaseUse(worker.Name)
		if lowestPort == 0 {
			return nil
		}

		job.Address = "localhost:" + strconv.Itoa(lowestPort)
		schedule(job, worker.Name)
	}
	return nil
}