package models

const (
	GENDER_UNSPECIFIED = iota
	GENDER_MALE
	GENDER_FEMALE
)

const (
	ROLE_ACTOR = iota
	ROLE_DIRECTOR
	ROLE_WRITTER
	ROLE_MUSICAL
)

type Profile struct {
	Id        *int64  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Age       *int32  `json:"age"`
	Gender    *int    `json:"gender"`
	City      *string `json:"city"`
	Role      *int    `json:"role"`
}
