package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Client represents the base Nexus API client.
type Client struct {
	BaseURL     string
	restyClient *resty.Client

	// Services
	RepositoryService *RepositoryService
}

// NewClient initializes and returns a new Client instance.
func NewClient(baseURL, username, password string) (*Client, error) {
	// Basic validation for required fields
	if baseURL == "" || username == "" || password == "" {
		return nil, fmt.Errorf("baseURL, username, and password must all be provided")
	}

	// Trim trailing slashes from baseURL
	baseURL = strings.TrimRight(baseURL, "/")

	// Validate URL format
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL format: %w", err)
	}

	// Initialize the Resty client
	restyClient := resty.New().
		SetBaseURL(baseURL).
		SetBasicAuth(username, password).
		SetHeader("Content-Type", "application/json")

	client := &Client{
		BaseURL:     baseURL,
		restyClient: restyClient,
	}

	// Initialize services
	client.RepositoryService = &RepositoryService{client: client}

	return client, nil
}
