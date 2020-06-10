package controller

import (
	"dc-labs/mangos/protocol/surveyor"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.nanomsg.org/mangos"
	//"go.nanomsg.org/mangos/protocol/pub"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

var controllerAddress = "tcp://localhost:40899"
var sock mangos.Socket

type Worker struct {
	Name     string `json:"name"`
	Tags     string `json:"tags"`
	Status   string `json:"status"`
	Usage    int    `json:"usage"`
	URL      string `json:"url"`
	Token    string `json:"token"`
	Active   bool   `json:"active"`
	Port     int    `json:"port"`
	JobsDone int    `json:"jobsDone"`
}
type Test struct {
	id     int
	worker string
}
type ImageService struct{
	id int
	image string
	worker string
}
type Workload struct{
	folder string
	LastIndex int
}

var Done = make(chan string)

//var tests []Test
var tests = make(map[string]Test)
var Nodes = make(map[string]Worker)
var filters = make(map[string]ImageService)
var Results = make(map[string]Workload)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Start() {
	var sock mangos.Socket
	var err error
	var msg []byte
	if sock, err = surveyor.NewSocket(); err != nil {
		die("can't get new surveyor socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on surveyor socket: %s", err.Error())
	}
	err = sock.SetOption(mangos.OptionSurveyTime, time.Second)
	if err != nil {
		die("SetOption(): %s", err.Error())
	}
	for {
		err = sock.Send([]byte("Hello workers"))
		if err != nil {
			die("No workers %+v", err.Error())
		}

		for {
			if msg, err = sock.Recv(); err != nil {
				break
			}
			isRegistered := false
			worker := ParseResponse(string(msg))

			for _, v := range Nodes {
				if v.Name == worker.Name {
					isRegistered = true
				}
			}

			if !isRegistered {
				Nodes[worker.Name] = worker
			}
			fmt.Println(Nodes[worker.Name].Name, " serves in localhost:", Nodes[worker.Name].Port, "\n")
			// Could also use sock.RecvMsg to get header
		}
	}

}
func ParseResponse(msg string) Worker {
	worker := Worker{}
	data := strings.Split(msg, " ")
	worker.Name = data[0]
	worker.Status = "free"
	usage, _ := strconv.Atoi(data[2])
	worker.Usage = usage
	worker.Tags = data[3]
	port, _ := strconv.Atoi(data[4])
	worker.Port = port
	jobsDone, _ := strconv.Atoi(data[5])
	worker.JobsDone = jobsDone
	worker.Active = true
	worker.URL = "localhost:" + data[4]
	worker.Token = data[6]
	return worker
}

func IncreaseUse(name string) {
	if thisProduct, ok := Nodes[name]; ok {
		thisProduct.Usage++
		thisProduct.JobsDone++
		Nodes[name] = thisProduct
	}
}
func ChangeStatus(name string) {
	if thisProduct, ok := Nodes[name]; ok {
		if thisProduct.Status == "free" {
			thisProduct.Status = "in use"
		} else {
			thisProduct.Status = "free"
		}
	}
}
func GetWorker(id int) string {
	name := tests[strconv.Itoa(id)].worker
	return name
}
func Register(name string, num int) {
	tests[strconv.Itoa(num)] = Test{id: num, worker: name}
}

func RegisterImage(name string, filepath string, workload string, num int) {
	isRegistered := false
	for t, _ := range Results {
		if t == workload {
			isRegistered = true
			break
		}
	}
	if(!isRegistered){
		curr := Workload{folder: filepath, LastIndex: 0}
		Results[workload] = curr
	}else{
		curr := Workload{folder: filepath, LastIndex: Results[workload].LastIndex + 1}
		Results[workload] = curr
	}
	filters[strconv.Itoa(num)] = ImageService{id: Results[workload].LastIndex , image: filepath, worker: name}
}

func IsTokenValid(workerToken string)(bool){

	for _, v := range Nodes {
		if v.Token == workerToken {
			return true
		}
	}
	return false
}
