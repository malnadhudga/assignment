package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <log_file>")
		return
	}

	logFilePath := os.Args[1]

	file, err := os.Open(logFilePath)

	if err != nil {
		fmt.Printf("Error while opening file %s: %v\n", logFilePath, err)
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Error closing file %s: %v\n", logFilePath, closeErr)
		}
	}()

	infoCount := 0
	warningCount := 0
	errorCount := 0
	totalLines := 0

	scanner := bufio.NewScanner(file)
	// Iterate over each line in the file
	for scanner.Scan() {
		line := scanner.Text()
		totalLines++

		// Use switch to categorize log levels
		switch {
		case strings.Contains(line, "[INFO]"):
			infoCount++
		case strings.Contains(line, "[WARNING]"):
			warningCount++
		case strings.Contains(line, "[ERROR]"):
			errorCount++
		}
	}

	// Print the summary report
	fmt.Printf("Log Analysis of file: %s\n\n", logFilePath)
	fmt.Printf("INFO: %d entries\n", infoCount)
	fmt.Printf("WARNING: %d entries\n", warningCount)
	fmt.Printf("ERROR: %d entries\n", errorCount)

	// Bonus: Count total lines and calculate percentages
	fmt.Printf("\nTotal log lines: %d\n", totalLines)
	if totalLines > 0 {
		fmt.Printf("INFO percentage: %.2f%%\n", float64(infoCount)/float64(totalLines)*100)
		fmt.Printf("WARNING percentage: %.2f%%\n", float64(warningCount)/float64(totalLines)*100)
		fmt.Printf("ERROR percentage: %.2f%%\n", float64(errorCount)/float64(totalLines)*100)
	}

	// Added time-based message
	fmt.Printf("\nAnalyzed at: %s\n", time.Now().Format("2025-06-12 15:04:05"))
}
