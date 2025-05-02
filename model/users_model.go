package model

type Users struct {
	UserID    int64  `gorm:"primaryKey;column:user_id;autoIncrement"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Username  string `gorm:"column:username;unique"`
	DOB       string `gorm:"column:dob"`
}

// TableName specifies the table name for the Users model
func (Users) TableName() string {
	return "users"
}
