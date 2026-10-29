package auth

import (
	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
)

// BuildDefaultTokenStore selects the store the normalised cache mode names, keyed by the credential fingerprint.
func BuildDefaultTokenStore(config wccfg.SdkConfig) (TokenStore, error) {
	credentials := config.Auth.Normalised()
	namespace := BuildCredentialsNamespace(credentials.BaseURL, credentials.ClientID, credentials.ClientSecret,
		credentials.AuthString, credentials.GrantType)
	switch credentials.TokenCacheMode {
	case wccfg.TokenCacheMemoryOnly:
		return NewMemoryTokenStore(), nil
	case wccfg.TokenCacheExplicitPath:
		return NewFileTokenStore(credentials.TokenCachePath, namespace)
	default: // wccfg.TokenCachePlatformDefault
		path, err := BuildDefaultTokenCachePath(DefaultTokenCacheFilename)
		if err != nil {
			return nil, err
		}
		return NewFileTokenStore(path, namespace)
	}
}
