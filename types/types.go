package types

import "time"

// Activity represents activity such as education or event
// swagger:model
type Activity struct {
	ID            int       `json:"id" exapmle:"1"`
	Name          string    `json:"name" exapmle:"Amber history"`
	Description   string    `json:"description" exapmle:"Education about amber"`
	BasePrice     float32   `json:"basePrice" exapmle:"20.50"`
	CreationDate  time.Time `json:"creationDate" exapmle:"2024-10-08 14:23:45.6789013 +0000UTC"`
	Hidden        bool      `json:"hidden" example:"false"`
	Verified      bool      `json:"verified" exapmle:"true"`
	Category      string    `json:"category" exapmle:"Education"`
	AverageRating float32   `json:"averageRating" exapmle:"3.5"`
	FkPackageID   int       `json:"fk_Packageid" exapmle:"1"`
}

// Activity represents package created by organizer which can be combined of many different activities
// swagger:model
type Package struct {
	ID            int     `json:"id" exapmle:"1"`
	Name          string  `json:"name" example:"Amber"`
	Description   string  `json:"description" exapmle:"All educations about amber"`
	Price         float32 `json:"price" exapmle:"100.20"`
	FkOrganizerID int     `json:"fk_Organizerid" exapmle:"1"`
}

type Location struct {
	ID           int     `json:"id" exapmle:"1"`
	Address      string  `json:"address" exapmle:"Kaunas city"`
	Longitude    float64 `json:"longitude" exapmle:"50.215458"`
	Latitude     float64 `json:"latitude" exapmle:"50.459414"`
	FkActivityID int     `json:"fk_Activityid" exapmle:"1"`
}

// User represents first authorized system role
// swagger:model
type User struct {
	ID               int       `json:"id" exapmle:"1"`
	Username         string    `json:"username" exapmle:"user"`
	Password         string    `json:"password" exapmle:"password"`
	Email            string    `json:"email" exapmle:"user@email.com"`
	RegistrationDate time.Time `json:"registrationDate" exapmle:"2024-10-08 14:23:45.6789013 +0000UTC"`
	LastLoginDate    time.Time `json:"lastLoginDate" exapmle:"2024-10-08 14:23:45.6789013 +0000UTC"`
}

// Review represents comments and ratings left in activity by other user
// swagger:model
type Review struct {
	ID           int       `json:"id" exapmle:"1"`
	Date         time.Time `json:"date" exapmle:"2024-10-08 14:23:45.6789013 +0000UTC"`
	Comment      *string   `json:"comment" exapmle:"Very nice education!"`
	Rating       int       `json:"rating" exapmle:"5"`
	FkUserID     int       `json:"fk_Userid" exapmle:"1"`
	FkActivityID int       `json:"fk_Activityid" exapmle:"1"`
}

// Administrator represents system role. Has rights on entire system
// swagger:model
type Administrator struct {
	ID            int `json:"id" exapmle:"1"`
	SecurityLevel int `json:"securityLevel" exapmle:"2"`
}

// Organizer represents system role. Can create new packages and activities
// swagger:model
type Organizer struct {
	ID          int     `json:"id" exapmle:"1"`
	Description *string `json:"description" exapmle:"Organizes educations about amber"`
}

type Subscribers struct {
	ID               int       `json:"id" exapmle:"1"`
	Email            string    `json:"email" exapmle:"subscriber@email.com"`
	SubscriptionDate time.Time `json:"subscriptionDate" exapmle:"2024-10-08 14:23:45.6789013 +0000UTC"`
}

type Image struct {
	ID          int       `json:"id" exapmle:"1"`
	Description *string   `json:"description" exapmle:"atl text"`
	FilePath    string    `json:"filePath" exapmle:"resources/images/"`
	Url         string    `json:"url" exapmle:"resources/images/url"`
	UploadTime  time.Time `json:"uploadTime" exapmle:"2024-10-08 14:23:45.6789013 +0000UTC"`
}

type EntityImage struct {
	ID         int    `json:"id" exapmle:"1"`
	EntityType string `json:"entityType" exapmle:"activity"`
	FkEntity   int    `json:"fk_entity" exapmle:"1"`
	FkImageID  int    `json:"fk_Imageid" example:"1"`
}

type Category string

const (
	CategoryEducation Category = "Education"
	CategoryEvent     Category = "Event"
	CategoryService   Category = "Service"
	CategoryOther     Category = "Other"
)

// Payloads

// UserPayload represents the payload for viewing user, creating a new user and updating information.
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

// TODO:
type CreateAdministratorPayload struct {
	ID            int `json:"id" validate:"required" example:"123"`
	SecurityLevel int `json:"securityLevel" validate:"required" example:"2"`
}

// LoginUserPayload represents the payload for logging in existing user.
// swagger:model
type LoginUserPayload struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// ActivityPayload represents the payload for creating activities and updating them.
// swagger:model
type ActivityPayload struct {
	Name        string  `json:"name" validate:"required" example:"Amber history"`
	Description string  `json:"description" validate:"required" example:"Educations about amber"`
	BasePrice   float32 `json:"basePrice" validate:"required" example:"15.50"`
	Hidden      bool    `json:"hidden" validate:"required" example:"1"`
	Category    string  `json:"category" validate:"required" example:"Education"`
	FkPackageID int     `json:"fk_Packageid" validate:"required" example:"1"`
}

// ActivityFilterPayload represents the payload for filtering activities.
// swagger:model
type ActivityFilterPayload struct {
	Name      string  `json:"name" example:"Amber history"`
	MinPrice  float32 `json:"minPrice" example:"15.50"`
	MaxPrice  float32 `json:"maxPrice" example:"15.50"`
	Category  string  `json:"category" example:"Education"`
	MinRating int     `json:"minRating" validate:"min=1,max=5" example:"1"`
	MaxRating int     `json:"maxRating" validate:"min=1,max=5" example:"5"`
	Organizer string  `json:"organizer" example:"user"`
}

// CreatePackagePayload represents the payload for creating packages.
// swagger:model
type CreatePackagePayload struct {
	Name          string  `json:"name" validate:"required" example:"Amber"`
	Description   string  `json:"description" validate:"required" example:"Everything about amber"`
	Price         float32 `json:"price" validate:"required" example:"40"`
	FkOrganizerID int     `json:"fk_Organizerid" validate:"required" example:"1"`
}

// ReviewPayload represents the payload for creating reviews and updating them.
// swagger:model
type ReviewPayload struct {
	Comment      string `json:"comment" validate:"required" example:"Very nice education"`
	Rating       int    `json:"rating" validate:"required,min=1,max=5" example:"5"`
	FkUserID     int    `json:"fk_Userid" validate:"required" example:"1"`
	FkActivityID int    `json:"fk_Activityid" validate:"required" example:"1"`
}

// Interfaces
type UserCastle interface {
	GetUserByID(id int) (*User, error)

	GetAdministratorByID(id int) (*Administrator, error)

	GetOrganizerByID(id int) (*Organizer, error)
	GetOrganizerByActivityID(activityID int) (*Organizer, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) error
	UpdateUser(User) error
	DeleteUser(id int) error
	ListUsers() ([]*User, error)

	CreateOrganizer(Organizer) error
	CreateAdministrator(Administrator) error
}

type ActivityCastle interface {
	CreateActivity(Activity) error
	GetActivityByID(id int) (*Activity, error)
	UpdateActivity(Activity) error
	DeleteActivity(id int) error
	GetActivityInsidePackageByName(activityName string, packageID int) (*Activity, error)
	ListActivities() ([]*Activity, error)
	FilterActivities(ActivityFilterPayload) ([]*Activity, error)

	GetPackageByID(id int) (*Package, error)
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
	ListReviewsFromPackage(id int) ([]*Review, error)
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
