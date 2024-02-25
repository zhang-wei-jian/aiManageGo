package models

type User struct {
	//username  string `gorm:"column:id" json:"id"`
	//passworld string `gorm:"column:name" json:"name"`

	Username string ` json:"username"`
	Password string ` json:"password"`
}
