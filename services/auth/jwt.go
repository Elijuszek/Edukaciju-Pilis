package auth

import (
	"educations-castle/configs"
	"educations-castle/types"
	"educations-castle/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/net/context"
)

type contextKey string

const UserKey contextKey = "userID"
const RoleKey contextKey = "role"

var tokenBlacklist = make(map[string]time.Time)

func WithJWTAuth(handlerFunc http.HandlerFunc, castle types.UserCastle, requiredRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		token, err := ValidateJWT(tokenString)

		if err != nil || !token.Valid {
			PermissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, err := strconv.Atoi(claims["userID"].(string))
		if err != nil {
			PermissionDenied(w)
			return
		}

		// Check if the user still exists
		user, err := castle.GetUserByID(userID)
		if err != nil || user == nil {
			PermissionDenied(w)
			return
		}

		// Retrieve role from claims
		role, ok := claims["role"].(string)
		if !ok {
			PermissionDenied(w)
			return
		}

		// Verify that user has the required role
		if !hasRequiredRole(role, requiredRoles) {
			PermissionDenied(w)
			return
		}

		// Add user ID and role to context
		ctx := context.WithValue(r.Context(), UserKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID int, role string) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"role":      role,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	if isTokenBlacklisted(tokenString) {
		return nil, fmt.Errorf("token is blacklisted")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})
}

func AddTokenToBlacklist(tokenString string, expiration int64) {
	expirationTime := time.Unix(expiration, 0)
	tokenBlacklist[tokenString] = expirationTime
}

func isTokenBlacklisted(tokenString string) bool {
	if expiration, ok := tokenBlacklist[tokenString]; ok {
		return time.Now().Before(expiration)
	}
	return false
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}

func hasRequiredRole(userRole string, requiredRoles []string) bool {
	// If no roles are required, all valid users should pass
	if len(requiredRoles) == 0 {
		return true
	}

	for _, role := range requiredRoles {
		if userRole == role {
			return true
		}
	}
	return false
}

func CreateRefreshToken(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.RefreshTokenExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := utils.GetTokenFromRequest(r)
	token, err := ValidateJWT(refreshToken)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, err := strconv.Atoi(claims["userID"].(string))
	if err != nil {
		http.Error(w, "Invalid user ID in refresh token", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	accessToken, err := CreateJWT([]byte(configs.Envs.JWTSecret), userID, claims["role"].(string))
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}

	// Respond with the new access token
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func CheckOwnership(r *http.Request, resourceOwnerID int) bool {
	role := r.Context().Value(RoleKey).(string)
	if role == "administrator" {
		return true // Administrators can modify any resource
	}

	userID := r.Context().Value(UserKey).(int)
	return userID == resourceOwnerID // Regular users can modify only their resources
}

// TODO: prevent running out of memory by cleaning up expired tokens
func cleanupExpiredTokens(blacklist map[string]time.Time) {
	currentTime := time.Now()
	for token, expiryTime := range blacklist {
		if currentTime.After(expiryTime) {
			delete(blacklist, token)
		}
	}
}
