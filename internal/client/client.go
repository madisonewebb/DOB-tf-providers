package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents the DevOps API client
type Client struct {
	HostURL    string
	HTTPClient *http.Client
}

// Engineer represents an individual engineer
type Engineer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// NewClient creates a new DevOps API client
func NewClient(host string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    host,
	}

	return &c, nil
}

// doRequest performs HTTP requests to the API
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

// GetEngineers retrieves all engineers
func (c *Client) GetEngineers() ([]Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var engineers []Engineer
	err = json.Unmarshal(body, &engineers)
	if err != nil {
		return nil, err
	}

	return engineers, nil
}

// GetEngineer retrieves a specific engineer by ID
// Since the API doesn't support individual engineer retrieval,
// we get all engineers and filter by ID
func (c *Client) GetEngineer(engineerID string) (*Engineer, error) {
	engineers, err := c.GetEngineers()
	if err != nil {
		return nil, err
	}

	for _, engineer := range engineers {
		if engineer.ID == engineerID {
			return &engineer, nil
		}
	}

	return nil, fmt.Errorf("engineer with ID %s not found", engineerID)
}

// CreateEngineer creates a new engineer
func (c *Client) CreateEngineer(engineer Engineer) (*Engineer, error) {
	rb, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var newEngineer Engineer
	err = json.Unmarshal(body, &newEngineer)
	if err != nil {
		return nil, err
	}

	return &newEngineer, nil
}

// UpdateEngineer updates an existing engineer
func (c *Client) UpdateEngineer(engineerID string, engineer Engineer) (*Engineer, error) {
	rb, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.HostURL, engineerID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var updatedEngineer Engineer
	err = json.Unmarshal(body, &updatedEngineer)
	if err != nil {
		return nil, err
	}

	return &updatedEngineer, nil
}

// DeleteEngineer deletes an engineer
func (c *Client) DeleteEngineer(engineerID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/engineers/%s", c.HostURL, engineerID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
