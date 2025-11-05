package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

type Client struct {
	restClient *api.RESTClient
	repo       repository.Repository
}

type Milestone struct {
	Number       int       `json:"number"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	DueOn        time.Time `json:"due_on"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	URL          string    `json:"html_url"`
}

type MilestoneInput struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
	DueOn       string `json:"due_on,omitempty"`
}

// NewClient creates a new API client
func NewClient(repoOverride string) (*Client, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}

	var repo repository.Repository
	if repoOverride != "" {
		repo, err = repository.Parse(repoOverride)
		if err != nil {
			return nil, fmt.Errorf("invalid repository format: %w", err)
		}
	} else {
		repo, err = repository.Current()
		if err != nil {
			return nil, fmt.Errorf("could not determine repository (use --repo flag): %w", err)
		}
	}

	return &Client{
		restClient: restClient,
		repo:       repo,
	}, nil
}

// ListMilestones lists all milestones in the repository
func (c *Client) ListMilestones(state string, sort string) ([]Milestone, error) {
	var milestones []Milestone

	path := fmt.Sprintf("repos/%s/%s/milestones", c.repo.Owner, c.repo.Name)

	// Add query parameters
	if state != "" {
		path += fmt.Sprintf("?state=%s", state)
	}
	if sort != "" {
		if state != "" {
			path += "&"
		} else {
			path += "?"
		}
		path += fmt.Sprintf("sort=%s", sort)
	}

	err := c.restClient.Get(path, &milestones)
	if err != nil {
		return nil, fmt.Errorf("failed to list milestones: %w", err)
	}

	return milestones, nil
}

// GetMilestone retrieves a specific milestone by number
func (c *Client) GetMilestone(number int) (*Milestone, error) {
	var milestone Milestone

	path := fmt.Sprintf("repos/%s/%s/milestones/%d", c.repo.Owner, c.repo.Name, number)

	err := c.restClient.Get(path, &milestone)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone: %w", err)
	}

	return &milestone, nil
}

// FindMilestoneByTitle finds a milestone by its title
func (c *Client) FindMilestoneByTitle(title string) (*Milestone, error) {
	milestones, err := c.ListMilestones("all", "")
	if err != nil {
		return nil, err
	}

	for _, m := range milestones {
		if m.Title == title {
			return &m, nil
		}
	}

	return nil, fmt.Errorf("milestone %q not found", title)
}

// CreateMilestone creates a new milestone
func (c *Client) CreateMilestone(input MilestoneInput) (*Milestone, error) {
	var milestone Milestone

	path := fmt.Sprintf("repos/%s/%s/milestones", c.repo.Owner, c.repo.Name)

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	err = c.restClient.Post(path, bytes.NewReader(body), &milestone)
	if err != nil {
		return nil, fmt.Errorf("failed to create milestone: %w", err)
	}

	return &milestone, nil
}

// UpdateMilestone updates an existing milestone
func (c *Client) UpdateMilestone(number int, input MilestoneInput) (*Milestone, error) {
	var milestone Milestone

	path := fmt.Sprintf("repos/%s/%s/milestones/%d", c.repo.Owner, c.repo.Name, number)

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	err = c.restClient.Patch(path, bytes.NewReader(body), &milestone)
	if err != nil {
		return nil, fmt.Errorf("failed to update milestone: %w", err)
	}

	return &milestone, nil
}

// DeleteMilestone deletes a milestone
func (c *Client) DeleteMilestone(number int) error {
	path := fmt.Sprintf("repos/%s/%s/milestones/%d", c.repo.Owner, c.repo.Name, number)

	err := c.restClient.Delete(path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	return nil
}

// ParseDueDate converts various date formats to ISO 8601
func ParseDueDate(dateStr string) (string, error) {
	if dateStr == "" {
		return "", nil
	}

	// Try to parse as various formats
	formats := []string{
		time.RFC3339,           // 2025-12-31T00:00:00Z
		"2006-01-02",           // 2025-12-31
		"2006-01-02T15:04:05",  // 2025-12-31T00:00:00
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t.UTC().Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("invalid date format (use YYYY-MM-DD or ISO 8601)")
}
