package videorea_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
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

type fakeRunner struct {
	args []string
	out  string
	err  error
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	f.args = append([]string{name}, args...)
	if f.err != nil {
		return nil, f.err
	}
	for i, a := range args {
		if a == "-o" && i+1 < len(args) {
			dest := filepath.Dir(args[i+1])
			_ = os.WriteFile(filepath.Join(dest, "audio.aac"), []byte("aac"), 0o644)
		}
	}
	return []byte(f.out), nil
}

func TestDownloadMediaAndMetadata(t *testing.T) {
	dest := t.TempDir()
	fr := &fakeRunner{out: `{"title":"T","description":"D","duration":61.5,"extractor_key":"Youtube"}`}
	in, err := videorea.DownloadWithRunner(context.Background(), fr, "https://youtu.be/abc", dest)
	if err != nil {
		t.Fatal(err)
	}
	if in.Title != "T" || in.Description != "D" || in.Duration != 61.5 || in.URL != "https://youtu.be/abc" {
		t.Errorf("metadados: %+v", in)
	}
	if in.Platform != videorea.PlatformYouTube || in.AudioPath != filepath.Join(dest, "audio.aac") {
		t.Errorf("plataforma/áudio: %+v", in)
	}
	joined := strings.Join(fr.args, " ")
	for _, flag := range []string{"yt-dlp", "--dump-json", "--no-playlist", "--extract-audio", "--audio-format aac", "--audio-quality 0"} {
		if !strings.Contains(joined, flag) {
			t.Errorf("flag ausente %q em %s", flag, joined)
		}
	}
}

func TestDownloadMediaAndMetadataFailures(t *testing.T) {
	ctx := context.Background()
	if _, err := videorea.DownloadWithRunner(ctx, &fakeRunner{err: errors.New("boom")}, "https://youtu.be/a", t.TempDir()); err == nil {
		t.Error("falha do yt-dlp deve propagar")
	}
	if _, err := videorea.DownloadWithRunner(ctx, &fakeRunner{out: "não é json"}, "https://youtu.be/a", t.TempDir()); err == nil {
		t.Error("saída inválida deve falhar")
	}
	if _, err := videorea.DownloadWithRunner(ctx, &fakeRunner{out: "{}"}, "ftp://x", t.TempDir()); err == nil {
		t.Error("URL inválida deve falhar")
	}
}
