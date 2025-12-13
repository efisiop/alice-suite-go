package auth

import (
	"testing"
)

// TestHashPassword tests password hashing
func TestHashPassword(t *testing.T) {
	password := "test-password-123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	if hash == "" {
		t.Fatal("HashPassword returned empty string")
	}
	
	if hash == password {
		t.Fatal("HashPassword returned plain text password")
	}
}

// TestVerifyPassword tests password verification
func TestVerifyPassword(t *testing.T) {
	password := "test-password-123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	
	// Test correct password
	err = VerifyPassword(hash, password)
	if err != nil {
		t.Fatalf("VerifyPassword failed for correct password: %v", err)
	}
	
	// Test incorrect password
	err = VerifyPassword(hash, "wrong-password")
	if err == nil {
		t.Fatal("VerifyPassword should fail for incorrect password")
	}
}

// TestHashPassword_DifferentHashes tests that same password produces different hashes
func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "test-password-123"
	
	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("HashPassword failed: %v", err1)
	}
	
	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("HashPassword failed: %v", err2)
	}
	
	// Hashes should be different (bcrypt includes salt)
	if hash1 == hash2 {
		t.Error("Same password should produce different hashes (bcrypt includes salt)")
	}
	
	// But both should verify correctly
	err1 = VerifyPassword(hash1, password)
	err2 = VerifyPassword(hash2, password)
	
	if err1 != nil || err2 != nil {
		t.Error("Both hashes should verify correctly")
	}
}

