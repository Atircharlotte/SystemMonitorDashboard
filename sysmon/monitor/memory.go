package monitor

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// MemState is just like a JS Object or TypeScript interface
type MemState struct {
	Total int
	Available int
}
func GetMemory() MemState {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// create an empty struct to hold data
	var stats MemState

	for scanner.Scan() {
		line := scanner.Text()

		// split the line into words
		fields := strings.Fields(line)

		// ensure having at least two fields to avoid crashing (?)
		if len(fields) < 2 {
			continue
		}
		// check the first word, parse the second word into an integer
		if fields[0] == "MemTotal:" {
			// convert string to an integer
			val, _ := strconv.Atoi(fields[1])
			stats.Total = val
		} else if fields[0] == "MemAvailable:" {
			val, _ := strconv.Atoi(fields[1])
			stats.Available = val
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading file: ", err)
	}

	return stats

}