package types

import "time"

type Activity struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	BasePrice     float32   `json:"basePrice"`
	CreationDate  time.Time `json:"creationDate"`
	Hidden        bool      `json:"hidden"`
	Verified      bool      `json:"verified"`
	Category      string    `json:"category"`
	AverageRating float32   `json:"averageRating"`
	Fk_PackageId  int       `json:"fk_Packageid"`
	Fk_ThemeId    int       `json:"fk_Themeid"`
}

type Package struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	Price          float32 `json:"price"`
	Fk_OrganizerId int     `json:"fk_Organizerid"`
}

type Theme struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Fk_OrganizerId int    `json:"fk_Organizerid"`
}

type Location struct {
	ID            int      `json:"id"`
	Address       string   `json:"address"`
	Longitude     *float64 `json:"longitude"`
	Latitude      *float64 `json:"latitude"`
	Fk_ActivityId int      `json:"fk_Activityid"`
}

// swagger:model
type User struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	Email            string    `json:"email"`
	RegistrationDate time.Time `json:"registrationDate"`
	LastLoginDate    time.Time `json:"lastLoginDate"`
}

type Review struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	Comment      *string   `json:"comment"`
	Rating       int       `json:"rating"`
	FkUserId     int       `json:"fk_Userid"`
	FkActivityId int       `json:"fk_Activityid"`
}

type Administrator struct {
	ID            int `json:"id"`
	SecurityLevel int `json:"securityLevel"`
}

type Organizer struct {
	ID          int     `json:"id"`
	Description *string `json:"description"`
}

type Subscribers struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	SubscriptionDate time.Time `json:"subscriptionDate"`
}

type Image struct {
	ID          int       `json:"id"`
	Description *string   `json:"description"`
	FilePath    string    `json:"filePath"`
	Url         string    `json:"url"`
	UploadTime  time.Time `json:"uploadTime"`
}

type EntityImage struct {
	ID         int    `json:"id"`
	EntityType string `json:"entityType"`
	FkEntity   int    `json:"fk_entity"`
	FkImageId  int    `json:"fk_Imageid"`
}

type Category string

const (
	CategoryEducation Category = "Education"
	CategoryEvent     Category = "Event"
	CategoryService   Category = "Service"
	CategoryOther     Category = "Other"
)

// Payloads

// RegisterUserPayload represents the payload for creating a new user.
// swagger:model
type RegisterUserPayload struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// Interfaces

type UserCastle interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

// Responses

// UserResponse represents the response structure for a user.
// swagger:model
type UserResponse struct {
	ID               int       `json:"id" example:"1"`
	Username         string    `json:"username" example:"john_doe"`
	Email            string    `json:"email" example:"john.doe@example.com"`
	RegistrationDate time.Time `json:"registrationDate" example:"2023-10-01T15:04:05Z07:00"`
	LastLoginDate    time.Time `json:"lastLoginDate" example:"2023-10-01T18:04:05Z07:00"`
}

// ErrorResponse represents an error response
// swagger:model
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid payload or user already exists"`
}
