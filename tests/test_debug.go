package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"user"`
}

func main() {
	baseURL := "http://localhost:8080"
	
	fmt.Println("=== LOGGING IN ===")
	loginReq := map[string]string{
		"email":    "user25@example.com",
		"password": "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	
	// Debug the login request
	fmt.Println("Login request payload:", string(loginBody))
	
	loginResp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		fmt.Println("Error logging in:", err)
		return
	}
	defer loginResp.Body.Close()

	loginResponseBody, _ := io.ReadAll(loginResp.Body)
	fmt.Println("Login response status:", loginResp.Status)
	fmt.Println("Login response body:", string(loginResponseBody))
	
	var loginResult LoginResponse
	if err := json.Unmarshal(loginResponseBody, &loginResult); err != nil {
		fmt.Println("Error parsing login response:", err)
		return
	}

	if loginResult.Token == "" {
		fmt.Println("No token received in login response!")
		return
	}

	fmt.Println("Got token:", loginResult.Token)
	fmt.Println("User ID:", loginResult.User.ID)

	fmt.Println("\n=== CHECKING TOKEN VALIDITY ===")
	// Let's first verify if our token is working by getting user info
	userReq, _ := http.NewRequest("GET", baseURL+"/api/users/me", nil)
	userReq.Header.Set("Authorization", "Bearer "+loginResult.Token)
	
	client := &http.Client{}
	userResp, err := client.Do(userReq)
	if err != nil {
		fmt.Println("Error getting user info:", err)
		return
	}
	defer userResp.Body.Close()
	
	userRespBody, _ := io.ReadAll(userResp.Body)
	fmt.Println("User info response status:", userResp.Status)
	fmt.Println("User info response body:", string(userRespBody))

	fmt.Println("\n=== CREATING BOOKING ===")

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 3)   
	
	bookingReq := map[string]interface{}{
		"room_id":   1, 
		"from_date": fromDate.Format(time.RFC3339),
		"to_date":   toDate.Format(time.RFC3339),
	}
	
	bookingBody, _ := json.MarshalIndent(bookingReq, "", "  ")
	fmt.Println("Booking request payload:", string(bookingBody))

	req, _ := http.NewRequest("POST", baseURL+"/api/bookings", bytes.NewBuffer(bookingBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResult.Token)
	
	// Dump the full request for debugging
	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println("Full request:")
	fmt.Println(string(reqDump))
	
	// Use a client with timeout to avoid hanging
	clientWithTimeout := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := clientWithTimeout.Do(req)
	if err != nil {
		fmt.Println("Error creating booking:", err)
		fmt.Println("Let's try with a different content type...")
		
		// Try with a different content type
		req, _ = http.NewRequest("POST", baseURL+"/api/bookings", bytes.NewBuffer(bookingBody))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")  // Try with explicit charset
		req.Header.Set("Authorization", "Bearer "+loginResult.Token)
		
		resp, err = clientWithTimeout.Do(req)
		if err != nil {
			fmt.Println("Still error with different content type:", err)
			return
		}
	}
	
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Booking response status:", resp.Status)
	fmt.Println("Booking response body:", string(respBody))
	
	fmt.Println("\n=== CHECKING USER BOOKINGS ===")
	bookingsReq, _ := http.NewRequest("GET", baseURL+"/api/bookings", nil)
	bookingsReq.Header.Set("Authorization", "Bearer "+loginResult.Token)
	
	bookingsResp, err := client.Do(bookingsReq)
	if err != nil {
		fmt.Println("Error getting bookings:", err)
		return
	}
	defer bookingsResp.Body.Close()
	
	bookingsRespBody, _ := io.ReadAll(bookingsResp.Body)
	fmt.Println("User bookings response status:", bookingsResp.Status)
	fmt.Println("User bookings response body:", string(bookingsRespBody))
}