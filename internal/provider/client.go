package provider

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

// Developer represents a collection of developer engineers
type Developer struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Engineers []Engineer `json:"engineers"`
}

// Operations represents a collection of operations engineers
type Operations struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Engineers []Engineer `json:"engineers"`
}

// DevOps represents a combination of developer and operations engineers
type DevOps struct {
	ID  string     `json:"id"`
	Dev Developer  `json:"dev"`
	Ops Operations `json:"ops"`
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
func (c *Client) GetEngineer(engineerID string) (*Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/%s", c.HostURL, engineerID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var engineer Engineer
	err = json.Unmarshal(body, &engineer)
	if err != nil {
		return nil, err
	}

	return &engineer, nil
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

// GetDevelopers retrieves all developer teams
func (c *Client) GetDevelopers() ([]Developer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/developers", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var developers []Developer
	err = json.Unmarshal(body, &developers)
	if err != nil {
		return nil, err
	}

	return developers, nil
}

// GetDeveloper retrieves a specific developer team by ID
func (c *Client) GetDeveloper(developerID string) (*Developer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/developers/%s", c.HostURL, developerID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var developer Developer
	err = json.Unmarshal(body, &developer)
	if err != nil {
		return nil, err
	}

	return &developer, nil
}

// GetOperations retrieves all operations teams
func (c *Client) GetOperations() ([]Operations, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/operations", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var operations []Operations
	err = json.Unmarshal(body, &operations)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

// GetOperation retrieves a specific operations team by ID
func (c *Client) GetOperation(operationID string) (*Operations, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/operations/%s", c.HostURL, operationID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var operation Operations
	err = json.Unmarshal(body, &operation)
	if err != nil {
		return nil, err
	}

	return &operation, nil
}

// GetDevOps retrieves all DevOps teams
func (c *Client) GetDevOps() ([]DevOps, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devops", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var devops []DevOps
	err = json.Unmarshal(body, &devops)
	if err != nil {
		return nil, err
	}

	return devops, nil
}

// GetDevOp retrieves a specific DevOps team by ID
func (c *Client) GetDevOp(devopID string) (*DevOps, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devops/%s", c.HostURL, devopID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var devop DevOps
	err = json.Unmarshal(body, &devop)
	if err != nil {
		return nil, err
	}

	return &devop, nil
}
