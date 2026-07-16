package authclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sv-blog/internal/platform/config"
)

var ErrUnauthorized = errors.New("authclient: unauthorized")

type User struct {
	ID             int64     `json:"id"`
	OrganizationID int64     `json:"organization_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	DisplayName    string    `json:"display_name"`
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Provider struct {
	Provider  string `json:"provider"`
	IsPrimary bool   `json:"is_primary"`
}

type App struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	PermissionCount int    `json:"permission_count"`
}

type MeData struct {
	User           User       `json:"user"`
	Providers      []Provider `json:"providers"`
	AppAccessCount int        `json:"app_access_count"`
	Apps           []App      `json:"apps"`
}

type meResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    MeData `json:"data"`
}

type AccessCheck struct {
	Allowed    bool   `json:"allowed"`
	App        string `json:"app"`
	Permission string `json:"permission"`
}

type checkResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    AccessCheck `json:"data"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.Auth.ServerURL, "/"),
		httpClient: &http.Client{
			Timeout: cfg.Auth.Timeout,
		},
	}
}

func (c *Client) Me(ctx context.Context, token string) (*MeData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/auth/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("authclient: unexpected status %d", resp.StatusCode)
	}

	var body meResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if !body.Success {
		return nil, ErrUnauthorized
	}

	return &body.Data, nil
}

func (c *Client) Check(ctx context.Context, token string, appCode string, permission string) (bool, error) {
	endpoint := c.baseURL + "/api/v1/auth/me/apps/" + url.PathEscape(appCode) + "/check?permission=" + url.QueryEscape(permission)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return false, ErrUnauthorized
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false, fmt.Errorf("authclient: unexpected status %d", resp.StatusCode)
	}

	var body checkResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, err
	}
	if !body.Success {
		return false, ErrUnauthorized
	}

	return body.Data.Allowed, nil
}
