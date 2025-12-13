package auth

import (
	"errors"

	"github.com/efisiopittau/alice-suite-go/internal/database"
	"github.com/efisiopittau/alice-suite-go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Register creates a new user account
func Register(email, password, firstName, lastName string) (*models.User, error) {
	// Check if user already exists
	existing, err := database.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	// Hash password
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         "reader",
		IsVerified:   false,
	}

	err = database.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

// Login authenticates a user and returns the user object
func Login(email, password string) (*models.User, error) {
	// Get user by email
	user, err := database.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	err = VerifyPassword(user.PasswordHash, password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID string) (string, error) {
	// Get user to get email and role
	user, err := database.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	return GenerateJWT(userID, user.Email, user.Role)
}

// VerifyToken verifies a JWT token and returns the user ID
func VerifyToken(token string) (string, error) {
	claims, err := ValidateJWT(token)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

// GetUserFromToken extracts user information from a JWT token
func GetUserFromToken(token string) (*models.User, error) {
	claims, err := ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	user, err := database.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}



