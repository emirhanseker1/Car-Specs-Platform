package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api"

func main() {
	// Wait for server to start
	time.Sleep(2 * time.Second)

	fmt.Println("=== Testing New API ===\n")

	// Test 1: Health check
	fmt.Println("1. Health Check")
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		fmt.Printf("   ❌ Failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ Status: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	// Test 2: Create Brand
	fmt.Println("\n2. Create Brand (BMW)")
	country := "Germany"
	brandData := map[string]interface{}{
		"name":    "BMW",
		"country": country,
	}
	brandJSON, _ := json.Marshal(brandData)
	resp, err = http.Post(baseURL+"/brands", "application/json", bytes.NewBuffer(brandJSON))
	if err != nil {
		fmt.Printf("   ❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		var brand map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&brand)
		fmt.Printf("   ✓ Created: ID=%v, Name=%s\n", brand["id"], brand["name"])
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   ❌ Status: %d, Body: %s\n", resp.StatusCode, string(body))
	}

	// Test 3: List Brands
	fmt.Println("\n3. List Brands")
	resp, err = http.Get(baseURL + "/brands")
	if err != nil {
		fmt.Printf("   ❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var brands []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&brands)
		fmt.Printf("   ✓ Found %d brand(s)\n", len(brands))
		for _, b := range brands {
			fmt.Printf("      - %s\n", b["name"])
		}
	} else {
		fmt.Printf("   ❌ Status: %d\n", resp.StatusCode)
	}

	// Test 4: Create Model
	fmt.Println("\n4. Create Model (3 Series)")
	bodyStyle := "Sedan"
	modelData := map[string]interface{}{
		"brand_id":   1,
		"name":       "3 Series",
		"body_style": bodyStyle,
	}
	modelJSON, _ := json.Marshal(modelData)
	resp, err = http.Post(baseURL+"/models", "application/json", bytes.NewBuffer(modelJSON))
	if err != nil {
		fmt.Printf("   ❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		var model map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&model)
		fmt.Printf("   ✓ Created: ID=%v, Name=%s\n", model["id"], model["name"])
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("   ❌ Status: %d, Body: %s\n", resp.StatusCode, string(body))
	}

	// Test 5: Search (empty)
	fmt.Println("\n5. Search Trims")
	resp, err = http.Get(baseURL + "/search")
	if err != nil {
		fmt.Printf("   ❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		results := result["results"].([]interface{})
		fmt.Printf("   ✓ Found %d trim(s)\n", len(results))
	} else {
		fmt.Printf("   ❌ Status: %d\n", resp.StatusCode)
	}

	fmt.Println("\n=== Tests Complete ===")
}
