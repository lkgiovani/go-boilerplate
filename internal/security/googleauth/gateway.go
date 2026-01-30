package googleauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
)

type GoogleGateway struct {
	client           *http.Client
	allowedClientIDs []string
}

func NewGoogleGateway(androidClientID, iosClientID string) *GoogleGateway {
	return &GoogleGateway{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		allowedClientIDs: []string{androidClientID, iosClientID},
	}
}

func (g *GoogleGateway) VerifyAndExtract(ctx context.Context, idToken string) (*auth.GoogleUserInfo, error) {

	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google token validation failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
		Aud     string `json:"aud"`
		Exp     string `json:"exp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	isValidAudience := false
	for _, clientID := range g.allowedClientIDs {
		if result.Aud == clientID {
			isValidAudience = true
			break
		}
	}

	if !isValidAudience {
		return nil, fmt.Errorf("invalid google token audience: %s", result.Aud)
	}

	return &auth.GoogleUserInfo{
		Email:      result.Email,
		Name:       result.Name,
		PictureURL: result.Picture,
	}, nil
}
