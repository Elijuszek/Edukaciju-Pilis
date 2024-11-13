package auth

import (
	"educations-castle/configs"
	"educations-castle/types"
	"educations-castle/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/net/context"
)

type contextKey string

const UserKey contextKey = "userID"

const RoleKey contextKey = "role"

const RefreshTokenExpiry = time.Hour * 24

func WithJWTAuth(handlerFunc http.HandlerFunc, castle types.UserCastle, requiredRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		token, err := ValidateJWT(tokenString)

		if err != nil || !token.Valid {
			log.Printf("Invalid or failed token validation: %v", err)
			PermissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, err := strconv.Atoi(claims["userID"].(string))
		if err != nil {
			log.Printf("Failed to convert UserID to int: %v", err)
			PermissionDenied(w)
			return
		}

		// Retrieve role from claims
		role, ok := claims["role"].(string)
		if !ok {
			log.Printf("Role missing in token claims")
			PermissionDenied(w)
			return
		}

		// Verify that user has the required role
		if !hasRequiredRole(role, requiredRoles) {
			log.Printf("User does not have required role (required %s): %s", requiredRoles, role)
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

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})
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

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}

func hasRequiredRole(userRole string, requiredRoles []string) bool {
	for _, role := range requiredRoles {
		if userRole == role {
			return true
		}
	}
	return false
}

func CreateRefreshToken(secret []byte, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(RefreshTokenExpiry).Unix(),
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
