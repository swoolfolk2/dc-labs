package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	pb "github.com/CodersSquad/dc-labs/challenges/final/proto"
	"go.nanomsg.org/mangos"
	"google.golang.org/grpc"

	// register transports
	"go.nanomsg.org/mangos/protocol/respondent"
	_ "go.nanomsg.org/mangos/transport/all"
)

var (
	defaultRPCPort = 50051
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedTaskServer
}

var (
	controllerAddress = ""
	workerName        = ""
	tags              = ""
	status            = ""
	token             = ""
	endpoint          = ""
	workDone          = 0
	usage             = 0
	port              = 0
	jobsDone          = 0
	filter            = ""
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func init() {
	flag.StringVar(&controllerAddress, "controller", "tcp://localhost:40899", "Controller-Address")
	flag.StringVar(&workerName, "node-name", "hard-worker", "Worker-Name")
	flag.StringVar(&tags, "tags", "gpu,superCPU,largeMemory", "TAG1,TAG2, TAG3,TAGn")
	flag.StringVar(&token, "image-store-token", "token", "Img-Store-Token")
	flag.StringVar(&endpoint, "image-store-endpoint", "url_endpoint", "Image-Store-Endpoint")
}
func (s *server) SayHello(ctx context.Context, in *pb.TestRequest) (*pb.TestReply, error) {
	switch in.Name {
	case "test":
		workDone++
		log.Printf("RPC [Worker] %+v: testing...", workerName)
		usage++
		status = "Running"
		usage--
		return &pb.TestReply{Message: "Here " + workerName + " testing..."}, nil
	default:
		workDone++
		log.Printf("[Worker] %+v: calling", workerName)
		usage++
		status = "Running"
		return &pb.TestReply{Message: "Hello " + workerName}, nil
	}
}
func (s *server) FilterImg(ctx context.Context, in *pb.ImgRequest) (*pb.ImgReply, error) {

	fmt.Printf("I will filter the following image: ")
	fmt.Printf(in.GetImg().Filepath + "\n")
	fmt.Printf("Usign the following filter: ")
	fmt.Printf(in.GetImg().Filter + "\n")
	filter = in.GetImg().Filter
	// download image from APIs endpoint

	DownloadFile(in.GetImg().Filepath, in.Img.Index, in.Img.Workload, filter)

	return &pb.ImgReply{Message: "The image was proccesed by " + workerName}, nil

}

func DownloadFile(url string, index int64, workload string, filter string) (err error) {

	// Create the file
	name := strings.SplitN(url, ".", 2)
	filename := "./" + workload + "/" + strconv.Itoa(int(index)) + "." + name[1]
	_ = os.MkdirAll(workload+"/", 0755)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get("http://localhost:8080/" + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("i will download this image " + resp.Status)
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	makeFilter(url, filename, name[1], filter, workload)

	return nil
}

func makeFilter(url string, filename string, ext string, filter string, workload string) {
	switch filter {
	case "bw":
		workDone++
		log.Printf("RPC [Worker] %+v: will do bw filtering...", workerName)
		usage++

		cmd := exec.Command("test.exe", ""+filename, "bw")
		cmd.Run()

		status = "Running"
		usage--

		break
	case "sepia":
		workDone++
		log.Printf("RPC [Worker] %+v: will do sepia filtering...", workerName)
		usage++

		cmd := exec.Command("test.exe", ""+filename, "sepia")
		cmd.Run()

		status = "Running"
		usage--
		break
	case "avatar":
		workDone++
		log.Printf("RPC [Worker] %+v: will do avatar filtering...", workerName)
		usage++

		cmd := exec.Command("test.exe", ""+filename, "avatar")
		cmd.Run()
		status = "Running"
		usage--
		break
	default:
		break
	}

}

func UploadFile(filename string, workload string) (err error) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	/*Initial Request Values*/
	requestBody := &bytes.Buffer{}
	multiPartWriter := multipart.NewWriter(requestBody)
	fmt.Println(filepath.Base(file.Name()))
	/* Insert File */
	part, err := multiPartWriter.CreateFormFile("data", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
	}

	tokenSender, err := multiPartWriter.CreateFormField("worker-token")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = tokenSender.Write([]byte(token))
	if err != nil {
		log.Fatalln(err)
	}

	multiPartWriter.Close()

	url := "http://localhost:8080/upload/" + workload
	// Create the file
	request, err := http.NewRequest("POST", url, requestBody)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", multiPartWriter.FormDataContentType())

	/*Create Request and Send*/
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("[SendPostRequest()] Could not upload file.")
	}

}

// joinCluster is meant to join the controller message-passing server
func joinCluster() {
	var sock mangos.Socket
	var err error
	var msg []byte
	if sock, err = respondent.NewSocket(); err != nil {
		die("can't get new respondent socket: %s", err.Error())
	}
	log.Printf("Connecting to controller on: %s", controllerAddress)
	if err = sock.Dial(controllerAddress); err != nil {
		die("can't dial on respondent socket: %s", err.Error())
	}
	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		data := workerName + " " + status + " " + strconv.Itoa(usage) + " " + tags + " " + strconv.Itoa(defaultRPCPort) + " " + strconv.Itoa(jobsDone) + " " + token
		if err = sock.Send([]byte(data)); err != nil {
			die("Cannot send: %s", err.Error())
		}

		log.Printf("Message-Passing: Worker(%s): Received %s\n", workerName, string(msg))
	}
}

func getAvailablePort() int {
	port := defaultRPCPort
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			port = port + 1
			continue
		}
		ln.Close()
		break
	}
	return port
}

func main() {
	flag.Parse()
	//jobsDone := 0

	// Subscribe to Controller
	go joinCluster()

	// Setup Worker RPC Server
	rpcPort := getAvailablePort()
	defaultRPCPort = rpcPort
	log.Printf("Starting RPC Service on localhost:%v", rpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTaskServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
