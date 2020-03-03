package main

import(
	"os"
	"strconv"
	"fmt"
	"time"

)

func pipeline(a,b chan int){
	value := <- a
	b <- value 
}
/*
func main(){

	file,_ := os.Create("report1.txt")
	a := make(chan int)
	b := make(chan int)
	counter := 0
	
	start := time.Now()
	for{
		
		counter = counter + 1 
		go pipeline(a,b)
		a = b
		b = make(chan int)
		now := time.Since(start)
		fmt.Print("\rgoroutine: " + strconv.Itoa(counter)+ " in " + now.String())
		
		file.WriteString("\rgoroutine: " + strconv.Itoa(counter) +" in " + now.String())
	}
*/
	
}