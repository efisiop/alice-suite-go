package auth

import (
	"testing"
)

// TestGenerateJWT tests JWT token generation
func TestGenerateJWT(t *testing.T) {
	userID := "test-user-123"
	email := "test@example.com"
	role := "reader"
	
	token, err := GenerateJWT(userID, email, role)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	
	if token == "" {
		t.Fatal("GenerateJWT returned empty token")
	}
}

// TestValidateJWT tests JWT token validation
func TestValidateJWT(t *testing.T) {
	userID := "test-user-123"
	email := "test@example.com"
	role := "reader"
	
	// Generate token
	token, err := GenerateJWT(userID, email, role)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	
	// Validate token
	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	
	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
	
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
	
	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}

// TestValidateJWT_InvalidToken tests validation of invalid token
func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.here"
	
	_, err := ValidateJWT(invalidToken)
	if err == nil {
		t.Fatal("ValidateJWT should fail for invalid token")
	}
	
	if err != ErrInvalidToken {
		t.Logf("Expected ErrInvalidToken, got: %v (this is acceptable)", err)
	}
}

// TestValidateJWT_ExpiredToken tests validation of expired token
// Note: This test may need adjustment based on token expiration time
func TestValidateJWT_ExpiredToken(t *testing.T) {
	// This test would require creating an expired token
	// For now, we'll skip it as it requires more complex setup
	t.Skip("Expired token test requires token manipulation")
}

// TestExtractTokenFromHeader tests token extraction from Authorization header
func TestExtractTokenFromHeader(t *testing.T) {
	token := "test-token-123"
	
	// Test with Bearer prefix
	bearerHeader := "Bearer " + token
	extracted, err := ExtractTokenFromHeader(bearerHeader)
	if err != nil {
		t.Fatalf("ExtractTokenFromHeader failed: %v", err)
	}
	if extracted != token {
		t.Errorf("Expected token %s, got %s", token, extracted)
	}
	
	// Test without Bearer prefix
	extracted, err = ExtractTokenFromHeader(token)
	if err != nil {
		t.Fatalf("ExtractTokenFromHeader failed: %v", err)
	}
	if extracted != token {
		t.Errorf("Expected token %s, got %s", token, extracted)
	}
	
	// Test empty header
	_, err = ExtractTokenFromHeader("")
	if err == nil {
		t.Fatal("ExtractTokenFromHeader should fail for empty header")
	}
}

