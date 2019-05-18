package entitys

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

// User is our user table structure.
type User struct {
	ID        int64 // auto-increment by-default by xorm
	//Version   int   `xorm:"version"`
	Salt      string
	Username  string `xorm:"unique"`
	Password  string `xorm:"varchar(200)"`
	Nickname  string
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (u *User) GeneratePassword(userPassword string) bool {
	u.Salt = time.Now().String()
	u.Password = u.md5Password(userPassword)
	return u.Salt != "" && u.Password != ""
}

func (u *User) UpdatePassword(userPassword string) bool {
	if u.Salt == "" || userPassword == "" {
		return false
	}
	u.Password = u.md5Password(userPassword)
	return true
}

func (u *User) ValidatePassword(userPassword string, hashedPassword string) bool {
	return u.md5Password(userPassword) == hashedPassword
}

func (u *User) md5Password(userPassword string) string {
	m5 := md5.New()
	m5.Write([]byte(userPassword))
	m5.Write([]byte(u.Salt))
	return hex.EncodeToString(m5.Sum(nil))
}
