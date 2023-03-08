package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Whois struct {
	Available bool `json:"available"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Get the GoDaddy API credentials from the .env file
	key := os.Getenv("GODADDY_KEY")
	secret := os.Getenv("GODADDY_SECRET")

	file, err := os.Open("domains.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Create separate files for available and taken domains
	availableFile, err := os.Create("available.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer availableFile.Close()

	takenFile, err := os.Create("taken.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer takenFile.Close()

	// Loop through each domain in the file and check availability
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())

		// Send a GET request to the GoDaddy API to check availability
		url := fmt.Sprintf("https://api.godaddy.com/v1/domains/available?domain=%s", domain)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer res.Body.Close()

		// Parse the response from the GoDaddy API
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		var whois Whois
		err = json.Unmarshal(body, &whois)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return
		}

		// Write the domain to the appropriate file based on availability
		if whois.Available {
			availableFile.WriteString(domain + "\n")
			fmt.Printf("✅ %s is available!\n", domain)
		} else {
			takenFile.WriteString(domain + "\n")
			//fmt.Printf("❌ %s is taken.\n", domain)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("Done!")

}
