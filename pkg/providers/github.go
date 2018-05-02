package providers

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
)

// Header constants
const (
	XHubSignature   = "X-Hub-Signature"
	XGitHubEvent    = "X-GitHub-Event"
	XGitHubDelivery = "X-GitHub-Delivery"
)

const (
	SignaturePrefix = "sha1="
	SignatureLength = 45
)

type GithubProvider struct {
	secret string
}

func NewGithubProvider(secret string) (*GithubProvider, error) {
	if len(strings.TrimSpace(secret)) == 0 {
		return nil, errors.New("Cannot create github provider with empty secret")
	}

	return &GithubProvider{
		secret: secret,
	}, nil
}

func (p *GithubProvider) GetHeaderKeys() []string {
	return []string{
		XGitHubDelivery,
		XGitHubEvent,
		XHubSignature,
	}
}

// TODO: Update implementation and tests

// Github Signature Validation:
func (p *GithubProvider) Validate(hook Hook) bool {

	githubSignature := hook.Headers[XHubSignature]
	if len(githubSignature) != SignatureLength ||
		!strings.HasPrefix(githubSignature, SignaturePrefix) {
		return false
	}

	return IsValidPayload(p.secret, githubSignature[len(SignaturePrefix):], hook.Payload)

}

// IsValidPayload checks if the github payload's hash fits with
// the hash computed by GitHub sent as a header
func IsValidPayload(secret, headerHash string, payload []byte) bool {
	hash := HashPayload(secret, payload)
	return hmac.Equal(
		[]byte(hash),
		[]byte(headerHash),
	)
}

// HashPayload computes the hash of payload's body according to the webhook's secret token
// see https://developer.github.com/webhooks/securing/#validating-payloads-from-github
// returning the hash as a hexadecimal string
func HashPayload(secret string, playloadBody []byte) string {
	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write(playloadBody)
	sum := hm.Sum(nil)
	return fmt.Sprintf("%x", sum)
}
