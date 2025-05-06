package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	loginResp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		fmt.Println("Error logging in:", err)
		return
	}
	defer loginResp.Body.Close()

	loginResponseBody, _ := io.ReadAll(loginResp.Body)
	fmt.Println("Login response status:", loginResp.Status)
	
	var loginResult LoginResponse
	if err := json.Unmarshal(loginResponseBody, &loginResult); err != nil {
		fmt.Println("Error parsing login response:", err)
		fmt.Println("Login response body:", string(loginResponseBody))
		return
	}

	fmt.Println("Got token:", loginResult.Token)

	fmt.Println("\n=== CREATING BOOKING ===")

	fromDate := time.Now().AddDate(0, 0, 1)
	toDate := time.Now().AddDate(0, 0, 3)   
	
	bookingReq := map[string]interface{}{
		"room_id":   1, 
		"from_date": fromDate.Format(time.RFC3339),
		"to_date":   toDate.Format(time.RFC3339),
	}
	
	bookingBody, _ := json.MarshalIndent(bookingReq, "", "  ")
	fmt.Println("Booking request:", string(bookingBody))

	req, _ := http.NewRequest("POST", baseURL+"/api/bookings", bytes.NewBuffer(bookingBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResult.Token)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating booking:", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Booking response status:", resp.Status)
	fmt.Println("Booking response body:", string(respBody))
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
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
}