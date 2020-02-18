package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)
/*
Struct Package to define store all the values from an installed package.
name: string  with the name of a package.
installed: string with the date that the package was installed.
update: string with the date of the last update, if it was updated.
updates:  int with a counter of how many times a package was updated.
removal: string with the date of when the package was removed, if it was removed.
*/
type  Package struct {

	name string
	installed string
	update string
	updates int
	removal string

}
/*
Function that checks if there was an error, if so, exit the program
e: error sent to check
*/
func check(e error){

	if(e != nil){
		fmt.Println("Error")
		os.Exit(-1)
	}

}
/*
Function to write a string in a Txt file
toWrite: string that is going to be written.
file: File in which the string will be written
*/
func writeToTxt(toWrite string,file *os.File){

	_,err:=file.WriteString(toWrite)
	check(err)

}

// Main function
func main() {

	//shows that the program started
	fmt.Println("Pacman Log Analyzer\n")

	//Verify that the text has at least 2 lines
	if len(os.Args) < 2 {
		fmt.Println("You must send at least one pacman log file to analize")
		fmt.Println("usage: ./pacman_log_analizer <logfile>")
		os.Exit(1)
	}

	//setting counter variables to keep track of the packages
	var installed, removed, upgraded, current int

	//reads the given file and checks if error
	file, err := os.Open(os.Args[1])
	check(err)

	//Create a map so that the packages are not doubled
	packages := make(map[string]*Package)
	defer file.Close()

	//scan the file to read line by line
	scanner := bufio.NewScanner(file)

	//for each line of the scanner
	for scanner.Scan(){

		//the variable for the line of the text
		line := strings.Split(scanner.Text()," ")
		var name string
		var date string
		ope := line[3]
		if(len(line) > 4){
			name = line[4]
			date = line[0][1:len(line[0])] +" "+ line[1][:len(line[1])-1]
		}

		//depending on the keyword, it decides what to do with the package
		switch ope{

			case "installed":
                       		installed = installed + 1
				packages[name] = &Package{name,date,"-",0,"-"}

            		case "removed":
				packages[name].removal = date
				upgraded = upgraded - packages[name].updates
                        	removed = removed + 1
				installed = installed - 1
                	case "upgraded":
				packages[name].update = date
				packages[name].updates = packages[name].updates + 1
        	                if(packages[name].updates == 1){
					upgraded = upgraded + 1
				}

			case "reinstalled":
				packages[name].installed = date
				packages[name].update = "-"
				packages[name].updates = 0
				packages[name].removal = "-"
				removed = removed - 1
			default:

		}
	}

	//checks for errors in the scanner
	check(scanner.Err())

	//Creates a new empty Txt file
	fileTxt, err := os.Create("packages_report.txt")
	check(err)

	current = installed - removed

	//template of the Txt file
	writeToTxt("Pacman Packages Report\n",fileTxt)
	writeToTxt("----------------------\n",fileTxt)
	writeToTxt("- Installed packages : "+strconv.Itoa(installed)+"\n",fileTxt)
	writeToTxt("- Removed packages   : "+strconv.Itoa(removed)+"\n",fileTxt)
	writeToTxt("- Upgraded packages  : "+strconv.Itoa(upgraded)+"\n",fileTxt)
	writeToTxt("- Current installed  : "+strconv.Itoa(current)+"\n\n",fileTxt)


	//for each package, write its attributes
	for _,p := range packages{


		writeToTxt("- Package Name        : "+p.name+"\n",fileTxt)
		writeToTxt("  - Install date      : "+p.installed+"\n",fileTxt)
		writeToTxt("  - Last update date  : "+p.update+"\n",fileTxt)
		writeToTxt("  - How many updates  : "+strconv.Itoa(p.updates)+"\n",fileTxt)
		writeToTxt("  - Removal date      : "+p.removal+"\n",fileTxt)

	}

	//checks for errors in the txt file
	check(fileTxt.Close())

	//shows that the programm finished
	fmt.Println("Program finished")
	return

}
