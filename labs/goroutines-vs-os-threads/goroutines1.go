package main

import (
	"strconv"
	"os"
)

func s(last ,next chan int){
	value := <-last
	next <- value
}


func main(){

	file,_ := os.Create("report.txt")
	in := make(chan int)
	a := in
	b := make(chan int)
	counter := 0
	for{
		go s(a,b)
		a = b
		b = make(chan int)
		counter = counter + 1
		file.WriteString("\r goroutine "+strconv.Itoa(counter))
	}
}
