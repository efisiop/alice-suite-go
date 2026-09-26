package handlers

import (
	"os"
	"strings"
	"testing"
)

func TestConsultantLoginSynchronizesTokenForServerNavigation(t *testing.T) {
	page, err := os.ReadFile("../templates/consultant/login.html")
	if err != nil {
		t.Fatalf("read consultant login template: %v", err)
	}

	if !strings.Contains(string(page), "setAuthToken(token);") {
		t.Fatal("consultant login must use setAuthToken so protected page links receive auth_token cookie")
	}
}
