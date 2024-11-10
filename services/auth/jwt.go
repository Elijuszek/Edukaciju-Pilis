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
