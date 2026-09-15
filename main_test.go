package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"ascii-art-web/handlers"
	"ascii-art-web/pkg/banner"
	"ascii-art-web/pkg/parser"
	"ascii-art-web/pkg/renderer"
)

func TestParser(t *testing.T) {
	t.Run("simple text, no error", func(t *testing.T) {
		matrix, err := parser.Parse("hi")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matrix) != 1 || string(matrix[0]) != "hi" {
			t.Errorf("got %v, want [[h i]]", matrix)
		}
	})

	t.Run("multiple lines split on \\n", func(t *testing.T) {
		matrix, err := parser.Parse("ab\ncd")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matrix) != 2 || string(matrix[0]) != "ab" || string(matrix[1]) != "cd" {
			t.Errorf("got %v, want [[a b] [c d]]", matrix)
		}
	})

	t.Run("empty string yields one empty row", func(t *testing.T) {
		matrix, err := parser.Parse("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matrix) != 1 || len(matrix[0]) != 0 {
			t.Errorf("got %v, want a single empty row", matrix)
		}
	})

	t.Run("trailing newline produces trailing empty row", func(t *testing.T) {
		matrix, err := parser.Parse("hi\n")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(matrix) != 2 || len(matrix[1]) != 0 {
			t.Errorf("got %v, want [[h i] []]", matrix)
		}
	})

	t.Run("boundary characters are accepted", func(t *testing.T) {
		// space (32) and '~' (126) are the inclusive bounds.
		if _, err := parser.Parse(" ~"); err != nil {
			t.Errorf("unexpected error for boundary chars: %v", err)
		}
	})

	t.Run("control character below range is rejected", func(t *testing.T) {
		if _, err := parser.Parse("hi\x01there"); err == nil {
			t.Error("expected error for control character, got nil")
		}
	})

	t.Run("DEL (127) is rejected", func(t *testing.T) {
		if _, err := parser.Parse(string(rune(127))); err == nil {
			t.Error("expected error for DEL character, got nil")
		}
	})

	t.Run("non-ASCII rune is rejected", func(t *testing.T) {
		if _, err := parser.Parse("café"); err == nil {
			t.Error("expected error for non-ASCII rune 'é', got nil")
		}
	})
}

func buildBanner(maxChar rune) string {
	var sb strings.Builder
	sb.WriteString("header\n")
	for c := rune(32); c <= maxChar; c++ {
		for row := 0; row < 8; row++ {
			fmt.Fprintf(&sb, "%c-row%d\n", c, row)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func TestBannerLoad(t *testing.T) {
	t.Run("missing file returns error", func(t *testing.T) {
		_, err := banner.Load(filepath.Join(t.TempDir(), "nope.txt"))
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})

	t.Run("valid file loads full glyph blocks", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "font.txt")
		if err := os.WriteFile(path, []byte(buildBanner('9')), 0o644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		m, err := banner.Load(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		glyph, ok := m[' ']
		if !ok || len(glyph) != 8 {
			t.Errorf("space glyph = %v, want 8 lines present", glyph)
		}
		if _, ok := m['9']; !ok {
			t.Error("expected '9' to be present in bannerMap")
		}
	})

	t.Run("truncated file stops loading early", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "short.txt")

		if err := os.WriteFile(path, []byte(buildBanner('#')), 0o644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		m, err := banner.Load(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := m['#']; !ok {
			t.Error("expected '#' to be loaded")
		}
		if _, ok := m['Z']; ok {
			t.Error("did not expect 'Z' to be loaded from a truncated file")
		}
	})
}

func TestRenderer(t *testing.T) {
	t.Run("renders a known glyph map", func(t *testing.T) {
		bannerMap := map[rune][]string{
			'A': {"11", "22", "33", "44", "55", "66", "77", "88"},
			'B': {"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh"},
		}
		matrix := [][]rune{[]rune("AB")}

		got := string(renderer.Render(matrix, bannerMap))
		want := "11aa\n22bb\n33cc\n44dd\n55ee\n66ff\n77gg\n88hh\n"
		if got != want {
			t.Errorf("got:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("blank row renders a single newline", func(t *testing.T) {
		matrix := [][]rune{{}, []rune("A")}
		bannerMap := map[rune][]string{
			'A': {"1", "2", "3", "4", "5", "6", "7", "8"},
		}

		got := string(renderer.Render(matrix, bannerMap))
		if !strings.HasPrefix(got, "\n") {
			t.Errorf("expected leading blank line, got %q", got)
		}
	})

	t.Run("unsupported character panics (documents a known gap)", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("known issue: Render panics on unmapped rune: %v", r)
			}
		}()
		bannerMap := map[rune][]string{} // 'Z' deliberately missing
		renderer.Render([][]rune{[]rune("Z")}, bannerMap)
		t.Error("expected a panic for an unmapped rune, but none occurred")
	})
}

func hasFixtures(t *testing.T) bool {
	t.Helper()
	if _, err := os.Stat("banners/standard.txt"); err != nil {
		t.Skip("banners/standard.txt not found relative to CWD; run `go test` from the repo root")
		return false
	}
	if _, err := os.Stat("frontend/index.html"); err != nil {
		t.Skip("frontend/index.html not found relative to CWD; run `go test` from the repo root")
		return false
	}
	return true
}

func TestInputHandler(t *testing.T) {
	hasFixtures(t)

	t.Run("wrong method -> 405", func(t *testing.T) {
		for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
			req := httptest.NewRequest(method, "/ascii-art", nil)
			w := httptest.NewRecorder()
			handlers.InputHandler(w, req)
			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("%s: got %d, want %d", method, w.Code, http.StatusMethodNotAllowed)
			}
		}
	})

	t.Run("empty body -> 200 (no error on empty input per assignment)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(""))
		w := httptest.NewRecorder()
		handlers.InputHandler(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("non-printable character -> 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader("hi\x01there"))
		w := httptest.NewRecorder()
		handlers.InputHandler(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing banner file -> 500", func(t *testing.T) {

		prevWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("getwd: %v", err)
		}
		if err := os.Chdir(t.TempDir()); err != nil {
			t.Fatalf("chdir: %v", err)
		}
		defer func() {
			if err := os.Chdir(prevWd); err != nil {
				t.Fatalf("restore chdir: %v", err)
			}
		}()

		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader("hi"))
		w := httptest.NewRecorder()
		handlers.InputHandler(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("valid request -> 200 with rendered body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader("hi"))
		w := httptest.NewRecorder()
		handlers.InputHandler(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
		}
		if w.Body.Len() == 0 {
			t.Error("expected non-empty rendered body")
		}
	})
}

func TestIndexHandler(t *testing.T) {
	hasFixtures(t)

	t.Run("unknown path -> 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
		w := httptest.NewRecorder()
		handlers.IndexHandler(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("got %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("wrong method -> 405", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		handlers.IndexHandler(w, req)
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})

	t.Run("GET / -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handlers.IndexHandler(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("got %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestFontHandlers(t *testing.T) {
	cases := []struct {
		name    string
		handler http.HandlerFunc
		want    string
	}{
		{"Standard", handlers.StandardHandler, "banners/standard.txt"},
		{"Shadow", handlers.ShadowHandler, "banners/shadow.txt"},
		{"Thinkertoy", handlers.ThinkertoyHandler, "banners/thinkertoy.txt"},
	}

	for _, tc := range cases {
		t.Run(tc.name+"/GET->200", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			w := httptest.NewRecorder()
			tc.handler(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("got %d, want %d", w.Code, http.StatusOK)
			}
			if got := w.Body.String(); got != tc.want {
				t.Errorf("got body %q, want %q", got, tc.want)
			}
		})

		t.Run(tc.name+"/DELETE->405", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/x", nil)
			w := httptest.NewRecorder()
			tc.handler(w, req)
			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("got %d, want %d", w.Code, http.StatusMethodNotAllowed)
			}
		})
	}

	t.Cleanup(func() {
		req := httptest.NewRequest(http.MethodGet, "/standard", nil)
		handlers.StandardHandler(httptest.NewRecorder(), req)
	})
}

func TestFontHandlers_ConcurrentAccess(t *testing.T) {

	var wg sync.WaitGroup
	calls := []http.HandlerFunc{handlers.StandardHandler, handlers.ShadowHandler, handlers.ThinkertoyHandler}

	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			calls[i%len(calls)](httptest.NewRecorder(), req)
		}(i)
	}
	wg.Wait()

	t.Cleanup(func() {
		req := httptest.NewRequest(http.MethodGet, "/standard", nil)
		handlers.StandardHandler(httptest.NewRecorder(), req)
	})
}
