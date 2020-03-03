package main

import(
	"os"
	"strconv"
	"fmt"
	"time"

)

func pong(a,b chan int){
	value := <- a
	b <- value
	
}

func main(){

	file,_ := os.Create("report2.txt")
	a := make(chan int)
	b := make(chan int)
	counter := 0
	start := time.Now()
	a <- counter
	
	for i:=0;i<2000000;i++{
		
		
		
		go pong(a,b)
		now := time.Since(start)
		fmt.Print("\r"+ strconv.Itoa(counter)+" pongs in "+ now.String())
		file.WriteString("\r"+strconv.Itoa(counter)+" pongs in "+ now.String())
		
		counter = <- b
		b <- counter + 1
		go pong(b,a)
		now = time.Since(start)
		fmt.Print("\r"+ strconv.Itoa(counter)+" pongs in "+ now.String())
		file.WriteString("\r"+strconv.Itoa(counter)+" pongs in "+ now.String())
		
		counter = <- a
		a <- counter + 1
		
		
	}
	
	
}