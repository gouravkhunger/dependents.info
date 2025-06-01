package github

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"

	"dependents.info/internal/config"
	"dependents.info/pkg/utils"
)

type OIDCService struct {
	issuer   string
	audience string
}

type GitHubClaims struct {
	Repository string `json:"repository"`
}

func NewOIDCService(config *config.Config) *OIDCService {
	return &OIDCService{
		issuer:   config.GitHubOIDCIssuer,
		audience: config.GitHubOIDCAudience,
	}
}

func (s *OIDCService) VerifyToken(ctx context.Context, rawToken string, expectedRepo string) error {
	provider, err := oidc.NewProvider(ctx, s.issuer)
	if err != nil {
		return fmt.Errorf("failed to get provider: %v", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: s.audience,
	})

	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		return fmt.Errorf("failed to verify token: %v", err)
	}

	var claims GitHubClaims

	if err := idToken.Claims(&claims); err != nil {
		return fmt.Errorf("failed to parse claims: %v", err)
	}

	if !utils.ValidateRepository(claims.Repository) {
		return fmt.Errorf("unexpected repository claim format: %s", claims.Repository)
	}

	if claims.Repository != expectedRepo {
		return fmt.Errorf("repository mismatch: token repo is %s, expected %s", claims.Repository, expectedRepo)
	}

	return nil
}
