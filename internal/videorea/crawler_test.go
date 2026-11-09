package videorea_test

import (
	"testing"

	"legoloop/internal/videorea"
)

func TestDetectPlatform(t *testing.T) {
	cases := []struct {
		url      string
		expected videorea.Platform
	}{
		{"https://www.youtube.com/watch?v=FJZkfjM8WhI", videorea.PlatformYouTube},
		{"https://youtu.be/FJZkfjM8WhI", videorea.PlatformYouTube},
		{"https://www.instagram.com/reel/C7xY9Z1A2bC/", videorea.PlatformInstagram},
		{"https://www.tiktok.com/@user/video/71234567890", videorea.PlatformTikTok},
		{"file://C:/videos/tutorial.mp4", videorea.PlatformLocal},
		{"local://scratch/test.mp4", videorea.PlatformLocal},
	}

	for _, c := range cases {
		p, err := videorea.DetectPlatform(c.url)
		if err != nil {
			t.Errorf("DetectPlatform(%q) retornou erro inesperado: %v", c.url, err)
		}
		if p != c.expected {
			t.Errorf("DetectPlatform(%q) = %q; esperava %q", c.url, p, c.expected)
		}
	}
}

func TestInvalidURLHandling(t *testing.T) {
	invalids := []string{
		"",
		"   ",
		"ftp://malformed-url.com",
		"https://unknownplatform.org/video/123",
	}

	for _, u := range invalids {
		_, err := videorea.DetectPlatform(u)
		if err != videorea.ErrInvalidURL {
			t.Errorf("DetectPlatform(%q) deveria falhar com ErrInvalidURL; obteve %v", u, err)
		}
	}
}
