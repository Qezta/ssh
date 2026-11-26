package tui

import "testing"

func TestIdentityFor(t *testing.T) {
	cases := map[string]identity{
		"hi":      identityQezta,
		"hello":   identityQezta,
		" HELLO ": identityQezta,
		"divit":   identityDivit,
		"div":     identityDivit,
		"qezta":   identityDivit,
		"guest":   identityDivit,
		"":        identityDivit,
	}
	for user, want := range cases {
		if got := identityFor(user); got != want {
			t.Errorf("identityFor(%q) = %v, want %v", user, got, want)
		}
	}
}
