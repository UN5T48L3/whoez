package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/likexian/whois"
)

func main() {
	// Open domains file
	file, err := os.Open("domains.txt")
	if err != nil {
		fmt.Println("Failed to open file: ", err)
		return
	}
	defer file.Close()

	// Create output files
	positiveFile, err := os.Create("available.txt")
	if err != nil {
		fmt.Println("Failed to create positive output file: ", err)
		return
	}
	defer positiveFile.Close()

	negativeFile, err := os.Create("taken.txt")
	if err != nil {
		fmt.Println("Failed to create negative output file: ", err)
		return
	}
	defer negativeFile.Close()

	// Scan domains file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if domain == "" {
			continue
		}

		// Check if domain is registered
		if isRegistered(domain) {
			fmt.Fprintln(negativeFile, domain)
		} else {
			fmt.Fprintln(positiveFile, domain)

			// Print unregistered domains to console in green color
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("Unregistered domain: %s\n", green(domain))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Failed to read file: ", err)
		return
	}

	fmt.Println("Done!")
}

// Check if domain is registered
func isRegistered(domain string) bool {
	raw, err := whois.Whois(domain)
	if err != nil {
		return false
	}
	if strings.Contains(raw, "No match for domain") || strings.Contains(raw, "NOT FOUND") || strings.Contains(raw, "Status: AVAILABLE") {
		return false
	}
	return true
}
