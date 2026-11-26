package tui

import (
	"strings"
	"testing"
)

func TestIdentityFor(t *testing.T) {
	cases := map[string]identity{
		"hi":      identityQezta,
		"hello":   identityQezta,
		" HELLO ": identityQezta,
		"divit":   identityDivit,
		"Div":     identityDivit,
		"qezta":   identityUnknown,
		"guest":   identityUnknown,
		"":        identityUnknown,
	}
	for user, want := range cases {
		if got := identityFor(user); got != want {
			t.Errorf("identityFor(%q) = %v, want %v", user, got, want)
		}
	}
}

func TestGuide(t *testing.T) {
	for _, u := range []string{"divit", "div", "hi", "hello"} {
		if g := Guide(u); g != "" {
			t.Errorf("Guide(%q) should be empty, got %q", u, g)
		}
	}
	g := Guide("qezta")
	if !strings.Contains(g, "ssh divit@qezta.com") || !strings.Contains(g, "ssh hi@qezta.com") {
		t.Errorf("Guide(qezta) missing pointers: %q", g)
	}
}
