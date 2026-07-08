package firebase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	firebaseSDK "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	authClient *auth.Client
	initOnce   sync.Once
	initErr    error
)

func InitClient(ctx context.Context) (*auth.Client, error) {
	initOnce.Do(func() {
		var app *firebaseSDK.App
		credentialsPath, credentialsJSON, err := resolveCredentialsSource()
		if err != nil {
			initErr = err
			return
		}

		if credentialsPath != "" {
			app, initErr = firebaseSDK.NewApp(ctx, nil, option.WithCredentialsFile(credentialsPath))
		} else if len(credentialsJSON) > 0 {
			app, initErr = firebaseSDK.NewApp(ctx, nil, option.WithCredentialsJSON(credentialsJSON))
		} else {
			app, initErr = firebaseSDK.NewApp(ctx, nil)
		}
		if initErr != nil {
			return
		}

		authClient, initErr = app.Auth(ctx)
	})

	if initErr != nil {
		return nil, fmt.Errorf("initialize firebase auth: %w", initErr)
	}
	if authClient == nil {
		return nil, fmt.Errorf("firebase auth client is not initialized")
	}

	return authClient, nil
}

func VerifyIDToken(ctx context.Context, token string) (*auth.Token, error) {
	client, err := InitClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.VerifyIDToken(ctx, token)
}

func resolveCredentialsSource() (string, []byte, error) {
	for _, key := range []string{"GOOGLE_CREDS", "GOOGLE_APPLICATION_CREDENTIALS"} {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			if _, err := os.Stat("GOOGLE_CREDS.json"); err == nil {
				value = "GOOGLE_CREDS.json"
			}
		}
		if value == "" {
			continue
		}

		if looksLikeJSON(value) {
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(value), &parsed); err != nil {
				return "", nil, fmt.Errorf("parse %s as JSON: %w", key, err)
			}
			return "", []byte(value), nil
		}

		path := filepath.Clean(value)
		if _, err := os.Stat(path); err != nil {
			return "", nil, fmt.Errorf("read credentials file %s: %w", path, err)
		}
		return path, nil, nil
	}

	return "", nil, nil
}

func looksLikeJSON(value string) bool {
	trimmed := strings.TrimSpace(value)
	return strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
}
