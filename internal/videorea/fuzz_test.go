package videorea_test

import (
	"testing"

	"legoloop/internal/videorea"
)

func FuzzCrawler(f *testing.F) {
	// Sementes conhecidas
	seeds := []string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ",
		"https://www.instagram.com/reel/C1234567890/",
		"https://www.tiktok.com/@user/video/1234567890123456789",
		"file://C:/videos/sample.mp4",
		"local://assets/demo.mp4",
		"",
		"   ",
		"http://",
		"not-a-valid-url",
		"https://evil.com/fake?url=youtube.com",
		"https://youtube.com.attacker.com/",
		"file:///etc/passwd",
		"https://\x00malformed.url",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// DetectPlatform nunca deve entrar em panic ou travar
		plat, err := videorea.DetectPlatform(input)
		if err == nil {
			if plat != videorea.PlatformYouTube &&
				plat != videorea.PlatformInstagram &&
				plat != videorea.PlatformTikTok &&
				plat != videorea.PlatformLocal {
				t.Fatalf("plataforma retornada desconhecida: %v", plat)
			}
		}
	})
}
