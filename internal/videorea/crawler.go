package videorea

import (
	"net/url"
	"strings"
)

// DetectPlatform identifica a plataforma de vídeo a partir da URL.
func DetectPlatform(rawURL string) (Platform, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", ErrInvalidURL
	}
	if strings.HasPrefix(trimmed, "file://") || strings.HasPrefix(trimmed, "local://") {
		return PlatformLocal, nil
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Host == "" {
		return "", ErrInvalidURL
	}
	host := strings.ToLower(parsed.Host)
	return matchHostToPlatform(host)
}

func matchHostToPlatform(host string) (Platform, error) {
	switch {
	case strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be"):
		return PlatformYouTube, nil
	case strings.Contains(host, "instagram.com"):
		return PlatformInstagram, nil
	case strings.Contains(host, "tiktok.com"):
		return PlatformTikTok, nil
	default:
		return "", ErrInvalidURL
	}
}
