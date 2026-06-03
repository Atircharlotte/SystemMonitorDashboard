package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("run getMemInfo...")
	memInfoLocation := "/proc/meminfo" 
	file, err := os.Open(memInfoLocation);

	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to open file: ", err)
		return
	}

	defer file.Close()

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println("Failed to read files:", err)
		return
	}
	
	fmt.Println("Content of meminfo:\n", string(data))

}