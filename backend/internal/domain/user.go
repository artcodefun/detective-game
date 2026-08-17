package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `bson:"_id"`
	AccessTokenHash string     `bson:"access_token_hash"`
	Device          DeviceInfo `bson:"device"`
	CreatedAt       time.Time  `bson:"created_at"`
}

type DevicePlatform string

const (
	DevicePlatformIOS     DevicePlatform = "ios"
	DevicePlatformAndroid DevicePlatform = "android"
)

type DeviceInfo struct {
	Platform     DevicePlatform `bson:"platform"`
	Manufacturer string         `bson:"manufacturer"`
	Model        string         `bson:"model"`
	OSVersion    string         `bson:"os_version"`
}

func (d DeviceInfo) IsValid() bool {
	return (d.Platform == DevicePlatformIOS || d.Platform == DevicePlatformAndroid) &&
		d.Manufacturer != "" && d.Model != "" && d.OSVersion != ""
}

func NewUser(accessTokenHash string, device DeviceInfo) User {
	return User{
		ID:              uuid.New(),
		AccessTokenHash: accessTokenHash,
		Device:          device,
		CreatedAt:       time.Now(),
	}
}
