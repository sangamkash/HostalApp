package Admin

type ExcessType int

const (
	Full ExcessType = iota
	ReadOnly
	ReadAndWrite
)

type AdminUserDetail struct {
	Username     string     `json:"username" bson:"username" validate:"required,min=3,max=20"`
	Email        string     `json:"email" bson:"email" validate:"required,email"`
	Password     string     `json:"password" bson:"password" validate:"required,min=8,max=64"`
	ExcessLevel  ExcessType `json:"excess_level" bson:"excess_level" validate:"required"`
	RefreshToken string     `json:"refresh_token" bson:"refresh_token"`
}

type AdminLogin struct {
	Username string `json:"username" bson:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" bson:"password" validate:"required,min=8,max=64"`
}
