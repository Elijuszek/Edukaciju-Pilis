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

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=3,max=130"`
	Email    string `json:"email" validate:"required,email"`
}

// Interfaces

type UserCastle interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}
