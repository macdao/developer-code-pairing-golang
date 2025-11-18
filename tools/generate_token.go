// +build tools

package main

import (
	"fmt"
	"os"
	"strconv"

	"order-service/internal/adapter/web"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tools/generate_token.go <userID>")
		fmt.Println("Example: go run tools/generate_token.go 1001")
		os.Exit(1)
	}

	userID, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid userID: %v\n", err)
		os.Exit(1)
	}

	token, err := web.GenerateToken(userID)
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated JWT token for userID %d:\n", userID)
	fmt.Println(token)
	fmt.Println("\nUse it in your API requests:")
	fmt.Printf("Authorization: Bearer %s\n", token)
}
