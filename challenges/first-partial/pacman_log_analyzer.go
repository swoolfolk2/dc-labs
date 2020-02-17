
package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
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

func writeToTxt(toWrite string,f *os.File){

	_,err:=f.WriteString(toWrite)
	check(err)

}

func main() {
	fmt.Println("Pacman Log Analyzer\n")

	if len(os.Args) < 2 {
		fmt.Println("You must send at least one pacman log file to analize")
		fmt.Println("usage: ./pacman_log_analizer <logfile>")
		os.Exit(1)
	}

	var installed, removed, upgraded, current int

	file, err := os.Open(os.Args[1])
	check(err)
	packages := make(map[string]*Package)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		line := strings.Split(scanner.Text()," ")

		ope := line[3]

		if(ope == "installed" || ope == "removed" || ope == "upgraded"){
			var name string = line[4]
			var date string = line[0][1:len(line[0])] +" "+ line[1][:len(line[1])-1]
			if (ope == "installed"){
                       		installed = installed + 1
				packages[name] = &Package{name,date,"-",0,"-"}

			}
            		if (ope == "removed"){
				packages[name].removal = date
				upgraded = upgraded - packages[name].updates
                        	removed = removed + 1
               		}
                	if (ope == "upgraded"){
				packages[name].update = date
				packages[name].updates = packages[name].updates + 1
        	                if(packages[name].updates == 1){
					upgraded = upgraded + 1
				}
	                }
			if (ope == "reinstalled"){
				packages[name].installed = date
				packages[name].update = "-"
				packages[name].updates = 0
				packages[name].removal = "-"
				removed = removed - 1

			}
		}
	}

	check(scanner.Err())
	f, err := os.Create("packages_report.txt")
	check(err)
	installed = len(packages)
	current = installed - removed

	writeToTxt("Pacman Packages Report\n",f)
	writeToTxt("----------------------\n",f)
	writeToTxt("- Installed packages : "+strconv.Itoa(installed)+"\n",f)
	writeToTxt("- Removed packages   : "+strconv.Itoa(removed)+"\n",f)
	writeToTxt("- Upgraded packages  : "+strconv.Itoa(upgraded)+"\n",f)
	writeToTxt("- Current installed  : "+strconv.Itoa(current)+"\n\n",f)


	for _,p := range packages{


		writeToTxt("- Package Name        : "+p.name+"\n",f)
		writeToTxt("  - Install date      : "+p.installed+"\n",f)
		writeToTxt("  - Last update date  : "+p.update+"\n",f)
		writeToTxt("  - How many updates  : "+strconv.Itoa(p.updates)+"\n",f)
		writeToTxt("  - Removal date      : "+p.removal+"\n",f)

	}
	check(f.Close())


}
