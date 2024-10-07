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
	FkPackageID   int       `json:"fk_Packageid"`
}

type Package struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float32 `json:"price"`
	FkOrganizerID int     `json:"fk_Organizerid"`
}

type Location struct {
	ID           int     `json:"id"`
	Address      string  `json:"address"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	FkActivityID int     `json:"fk_Activityid"`
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
	FkUserID     int       `json:"fk_Userid"`
	FkActivityID int       `json:"fk_Activityid"`
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
	FkImageID  int    `json:"fk_Imageid"`
}

type Category string

const (
	CategoryEducation Category = "Education"
	CategoryEvent     Category = "Event"
	CategoryService   Category = "Service"
	CategoryOther     Category = "Other"
)

// Payloads

// UserPayload represents the payload for creating a new user and updating information.
// swagger:model
type UserPayload struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required,min=5,max=64" example:"password123"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

type CreateOrganizerPayload struct {
	ID          int    `json:"id" validate:"required" example:"123"`
	Description string `json:"description" validate:"required" example:"organizer"`
}

// LoginUserPayload represents the payload for logging in existing user.
// swagger:model
type LoginUserPayload struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// CreateActivityPayload represents the payload for creating activities.
// swagger:model
type CreateActivityPayload struct {
	Name        string  `json:"name" validate:"required" example:"Amber history"`
	Description string  `json:"description" validate:"required" example:"Educations about amber"`
	BasePrice   float32 `json:"basePrice" validate:"required" example:"15.50"`
	Hidden      bool    `json:"hidden" validate:"required" example:"1"`
	Category    string  `json:"category" validate:"required" example:"Education"`
	FkPackageID int     `json:"fk_Packageid" validate:"required" example:"1"`
}

// DeleteActivityPayload represents the payload for creating activities.
// swagger:model
type DeleteActivityPayload struct {
	Name        string `json:"name" validate:"required"`
	FkPackageID int    `json:"fk_PackageId" validate:"required"`
}

// CreatePackagePayload represents the payload for creating packages.
// swagger:model
type CreatePackagePayload struct {
	Name          string  `json:"name" validate:"required" example:"Amber"`
	Description   string  `json:"description" validate:"required" example:"Everything about amber"`
	Price         float32 `json:"price" validate:"required" example:"40"`
	FkOrganizerID int     `json:"fk_Organizerid" validate:"required" example:"1"`
}

// CreateReviewPayload represents the payload for creating reviews.
// swagger:model
type CreateReviewPayload struct {
	Comment      string `json:"comment" validate:"required" example:"Very nice education"`
	Rating       int    `json:"rating" validate:"required,min=1,max=5" example:"5"`
	FkUserID     int    `json:"fk_Userid" validate:"required" example:"1"`
	FkActivityID int    `json:"fk_Activityid" validate:"required" example:"1"`
}

// UpdateReviewPayload represents the payload for updating reviews.
// swagger:model
type UpdateReviewPayload struct {
	Comment string `json:"comment" validate:"required"`
	Rating  int    `json:"rating" validate:"required,gte=0,lte=5"` // Example validation for rating
}

// Interfaces
type UserCastle interface {
	GetUserByID(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) error
	UpdateUser(User) error
	DeleteUser(id int) error
	ListUsers() ([]*User, error)

	CreateOrganizer(Organizer) error
}

type ActivityCastle interface {
	CreateActivity(Activity) error
	GetActivityByID(id int) (*Activity, error)
	UpdateActivity(Activity) error
	DeleteActivity(id int) error
	GetActivityInsidePackageByName(activityName string, packageID int) (*Activity, error)

	CreatePackage(Package) error
	DeletePackage(id int) error
	GetPackageByName(name string) (*Package, error)
}

type ReviewCastle interface {
	CreateReview(Review) error
	GetReviewByID(id int) (*Review, error)
	UpdateReview(Review) error
	DeleteReviewByID(id int) error
	ListReviews() ([]*Review, error)

	GetReviewFromActivityByID(idActivity int, idUser int) (*Review, error)
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
