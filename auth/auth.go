package auth

import (
	"errors"

	"code.google.com/p/go.crypto/bcrypt"
)

var (
	roles map[string][]string

	ErrInvalidUserName = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
)

func init() {
	roles = make(map[string][]string, 1)
}

// RegisterRole registers a new role and its capabilities
func RegisterRole(role string, capabilities ...string) {
	roles[role] = capabilities
}

// Auth handles user authentication and authorization
type Auth struct {
	Username string
	Password []byte
	Roles    []string
}

// New returns a Auth object
func New() *Auth {
	return &Auth{}
}

// GeneratePassword generates a hashed password for the user account, using BCrypt algorithm
func (a *Auth) GeneratePassword(password []byte) error {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = hash
	return nil
}

// SetUsername sets the username
func (a *Auth) SetUsername(username string) {
	a.Username = username
}

// Check, checks the username / password pair to validate a user identity
func (a *Auth) Check(username string, password []byte) error {
	if username != a.Username {
		return ErrInvalidUserName
	}
	if err := a.checkPassword(password); err != nil {
		return ErrInvalidPassword
	}
	return nil
}

func (a *Auth) checkPassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(a.Password, password)
}

// AddRole adds the given role string to the account
func (a *Auth) AddRole(r string) {
	a.Roles = append(a.Roles, r)
}

// Can checks if the user has ANY of the given capabilities
func (a *Auth) Can(capabilities ...string) bool {
	for _, c := range capabilities {
		if a.can(c) {
			return true
		}
	}
	return false
}

// CanAll checks if the user has ALL of the given capabilities
func (a *Auth) CanAll(capabilities ...string) bool {
	for _, c := range capabilities {
		if !a.can(c) {
			return false
		}
	}
	return true
}

func (a *Auth) can(capability string) bool {
	for _, role := range a.Roles {
		for _, c := range roleCapabilities(role) {
			if c == capability {
				return true
			}
		}
	}
	return false
}

func roleCapabilities(role string) []string {
	if roles[role] != nil {
		return roles[role]
	}
	return []string{}
}