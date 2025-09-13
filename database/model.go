package database

import (
	"time"
	"gorm.io/gorm"
)
// User モデル
type User struct {
    gorm.Model
    Name     string `json:"name" gorm:"size:255;not null;unique"`
    Password string `json:"-" gorm:"size:255;not null"`

    Figures []Figure `json:"figures"` // ユーザーが持つ記録
}

// Figure モデル（ユーザーごとのカロリー記録）
type Figure struct {
    gorm.Model
    UserID  uint `json:"user_id"` // 外部キー
    User    User `json:"-" gorm:"foreignKey:UserID"`

    Calorie int `json:"calorie" gorm:"default:0"`

    RecordedAt time.Time `json:"recorded_at"`
}