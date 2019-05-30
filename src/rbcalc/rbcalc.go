package main

import (
"fmt"
"os"
)


func printUsage(){
	fmt.Println("Usage:\n-u/update update latest data into database\n")
	fmt.Println("-s/show show every ball hit info\n")
	fmt.Println("-g/get [param] get suggestions by different way\n")
}

func doUpdate(){
}

func doShowAll(){
}

func getSuggest(){
}

func main(){
	argc:=len(os.Args)
	if argc <2 {
		printUsage()
	}else{
		switch os.Args[1]{
		case "-u":
			fallthrough
		case "update":
			doUpdate()
		case "-s":
			fallthrough
		case "show":
			doShowAll()
		case "-g":
			fallthrough
		case "get":
				getSuggest()
		}
	}
}
