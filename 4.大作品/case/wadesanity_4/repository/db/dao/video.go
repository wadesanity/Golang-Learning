package dao

import (
	"videoGo/case/wadesanity_4/repository/db/model"

	"gorm.io/gorm"
)

type Mo interface {
	model.User | model.Video
}

func GetTotalByOpts[M Mo](d *gorm.DB, m M, opts ...Option) (int64, error) {
	var count int64
	db1 := d.Model(new(M))
	for _, opt := range opts {
		db1 = opt(db1)
	}
	db1 = db1.Count(&count)
	return count, db1.Error
}

func GetTakeByOpts[M Mo](d *gorm.DB, m *M, opts ...Option) (*M, error) {
	db1 := d.Model(new(M))
	for _, opt := range opts {
		db1 = opt(db1)
	}
	res := db1.Take(&m)
	return m, res.Error
}

func GetFindByOpts[M Mo](d *gorm.DB, m *M, opts ...Option) (*M, error) {
	db1 := d.Model(new(M))
	for _, opt := range opts {
		db1 = opt(db1)
	}
	res := db1.Find(&m)
	return m, res.Error
}

func GetRowUintByOpts[M Mo](d *gorm.DB, m *M, opts ...Option) (uint, error) {
	db1 := d.Model(new(M))
	for _, opt := range opts {
		db1 = opt(db1)
	}
	var count uint
	res := db1.Row().Scan(&count)
	return count, res
}

func ChangeByModel[M Mo](d *gorm.DB, m *M) (*M, error) {
	return m, d.Updates(m).Error
}
