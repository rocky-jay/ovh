package updater

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestReleaseSourceDefaultsToThisRepository(t *testing.T) {
	t.Setenv("OVH_UPDATE_API", "")
	if got := LatestReleaseURL(); got != "https://api.github.com/repos/rocky-jay/ovh/releases/latest" {
		t.Fatalf("unexpected default release source: %s", got)
	}
}

func TestFetchLatestOnlyAcceptsStableRelease(t *testing.T) {
	for _, test := range []struct {
		name    string
		release Release
		accept  bool
	}{
		{"stable", Release{TagName: "v0.4.38"}, true},
		{"draft", Release{TagName: "v0.4.39", Draft: true}, false},
		{"prerelease", Release{TagName: "v0.4.40-rc.1", Prerelease: true}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected release method: %s", r.Method)
				}
				if err := json.NewEncoder(w).Encode(test.release); err != nil {
					t.Error(err)
				}
			}))
			defer srv.Close()
			t.Setenv("OVH_UPDATE_API", srv.URL)
			release, err := FetchLatest()
			if test.accept {
				if err != nil || release == nil || release.TagName != test.release.TagName {
					t.Fatalf("stable release rejected: %+v, %v", release, err)
				}
			} else if err == nil || release != nil {
				t.Fatalf("non-public-stable release accepted: %+v, %v", release, err)
			}
		})
	}
}

func TestFetchLatestRejectsMalformedURL(t *testing.T) {
	t.Setenv("OVH_UPDATE_API", "http://[invalid")
	if _, err := FetchLatest(); err == nil {
		t.Fatal("malformed URL accepted")
	}
}

func TestPrepareRejectsMalformedChecksums(t *testing.T) {
	for _, sum := range []string{"bad", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123"} {
		t.Run(sum[:3], func(t *testing.T) {
			srv, rel := fakeRelease(t, []byte("new binary"), sum)
			defer srv.Close()
			exe := withFakeExe(t, "old binary")
			if _, _, err := Prepare(rel, nil); err == nil {
				t.Fatal("invalid SHA-256 checksum accepted")
			}
			body, err := os.ReadFile(exe)
			if err != nil || string(body) != "old binary" {
				t.Fatal("existing executable changed", err)
			}
		})
	}
}

func TestPrepareDoesNotReusePredictableTemporaryFile(t *testing.T) {
	srv, rel := fakeRelease(t, []byte("new binary"), "")
	defer srv.Close()
	exe := withFakeExe(t, "old binary")
	oldPath := filepath.Join(filepath.Dir(exe), ".ovh-server-update-probe")
	if err := os.WriteFile(oldPath, []byte("unrelated file"), 0600); err != nil {
		t.Fatal(err)
	}
	tmp, _, err := Prepare(rel, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)
	if body, err := os.ReadFile(oldPath); err != nil || string(body) != "unrelated file" {
		t.Fatal("update deleted an existing file while probing directory permissions", err)
	}
}
