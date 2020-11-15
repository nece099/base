package app

import (
	"time"

	"github.com/jinzhu/gorm"
)

const ( //
	SERVICE_TYPE_ALL       = "*"
	SERVICE_TYPE_UID       = "uid-service"
	SERVICE_TYPE_CONFIGURE = "configure-service"
	SERVICE_TYPE_AUTH      = "auth-service"
)

type TModel struct {
	ID        int64 `gorm:"AUTO_INCREMENT;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExternalService struct {
	AppID     string
	AppSecret string
	Address   string
}

type APP struct {
	ID          string `gorm:"primary_key;size:32"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `gorm:"size:64"` // 应用名称
	Remark      string `gorm:"size:64"` // 备注
	Secret      string `gorm:"size:64"`
	ServiceType string `gorm:"size:64"`
	Disabled    bool   `gorm:"default:0"`
}

func IsValidService(serviceType string) bool {
	if serviceType == SERVICE_TYPE_ALL ||
		serviceType == SERVICE_TYPE_UID ||
		serviceType == SERVICE_TYPE_CONFIGURE ||
		serviceType == SERVICE_TYPE_AUTH {
		return true
	}

	return false
}

func VerifyAppID(db *gorm.DB, appID string) error {
	app := &APP{}
	app.ID = appID

	err := db.Find(app).Error
	if err != nil {
		return err
	}
	return nil
}
