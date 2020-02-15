
package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"container/list"
)

type  Package struct {

	name string
	installed string 
	update string
	updates int
	removal string

}

func check(e error){

	if(e != nil){
		fmt.Println("Error")
		os.Exit(-1)
	}

}


func main() {
	fmt.Println("Pacman Log Analyzer")

	if len(os.Args) < 2 {
		fmt.Println("You must send at least one pacman log file to analize")
		fmt.Println("usage: ./pacman_log_analizer <logfile>")
		os.Exit(1)
	}

	var installed, removed, upgraded, current int

	file, err := os.Open(os.Args[1])
	check(err)
	packages := list.New()

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		line := strings.Split(scanner.Text()," ")

		ope := line[3]

		if(ope == "installed" || ope == "removed" || ope == "upgraded"){
			if (ope == "installed"){
                       		installed = installed + 1
                		packages.PushBack(Package{"prueba","02/02","02/04",2,""})

			}
            		if (ope == "removed"){
                        	removed = removed + 1
               		}
                	if (ope == "upgraded"){
        	                upgraded = upgraded + 1
	                }



		}
	}
	check(scanner.Err())


	fmt.Println(installed, removed, upgraded, current)


}
