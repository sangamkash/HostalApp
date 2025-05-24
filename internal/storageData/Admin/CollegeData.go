package Admin

type CollegeData struct {
	CollageName       string `json:"collage_name" bson:"collage_name" validate:"required,min=3,max=20"`
	CollageUniqueName string `json:"collage_unique_name" bson:"collage_unique_name" validate:"required,min=3,max=20"`
	CollageAddress    string `json:"collage_address" bson:"collage_address" validate:"required,min=3,max=20"`
	PinCode           string `json:"pin_code" bson:"pin_code" validate:"required,min=3,max=20"`
	CollageIcon       string `json:"collage_icon" bson:"collage_icon" validate:"required,min=3,max=20"`
	CollageStrength   int64  `json:"collage_strength" bson:"collage_strength" validate:"required,min=1,max=20"`
	MarkAsDeleted     bool   `json:"mark_as_deleted" bson:"mark_as_deleted" default:"false"`
}

type DelCollegeData struct {
	CollageUniqueName string `json:"collage_unique_name" bson:"collage_unique_name" validate:"required,min=3,max=20"`
}

type GetCollegeFilter struct {
	PinCode       string `json:"pin_code" bson:"pin_code"`
	MarkAsDeleted bool   `json:"mark_as_deleted" bson:"mark_as_deleted" default:"false"`
}
