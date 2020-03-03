package main

import (

	"os"
)

func giveValue(a,b chan int){
	v := <- a
	b <- v
}

func main(){
	f,_ := os.Create("report2.txt")
	a := make(chan int)
	b := make(chan int)
	value := 33
	for{
		go giveValue(a,b)
		go giveValue(b,a)
	}

	a <- value
}
