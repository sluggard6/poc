package model

import (
	"fmt"
	"time"

	"database/sql/driver"

	"github.com/sluggard/poc/util"
	"gorm.io/gorm"
)

//JsonTime time.Time的变体，为满足自定义序列化
type JsonTime time.Time

//MarshalJSON 实现它的json序列化方法
func (jtime JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(jtime).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (jtime JsonTime) Value() (driver.Value, error) {
	return time.Time(jtime), nil
}

//Scan 反射时转换为指针类型
func (jtime *JsonTime) Scan(src interface{}) error {
	*jtime = JsonTime(src.(time.Time))
	return nil
}

func (jtime *JsonTime) UnmarshalJSON(b []byte) error {
	// var err error
	t, err := time.Parse("2006-01-02 15:04:05", util.ByteArrayToString(b))
	*jtime = JsonTime(t)
	return err
}

//Model 替换gorm.Model
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt JsonTime       `json:"createdAt"`
	UpdatedAt JsonTime       `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleteAt"`
}

func DB() *gorm.DB {
	return db
}

func Create(model interface{}) (int64, error) {
	result := db.Create(model)
	return result.RowsAffected, result.Error
}

func GetById(model interface{}, id uint) (interface{}, error) {
	// result := db.First(model, id)
	result := db.Where("id = ?", id).First(model)
	return model, result.Error
}

func UpdateById(model interface{}) (int64, error) {
	result := db.Model(model).Updates(model)
	return result.RowsAffected, result.Error
}

func Delete(model interface{}) (int64, error) {
	result := db.Delete(model)
	return result.RowsAffected, result.Error
}
