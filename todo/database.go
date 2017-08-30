package main

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/lukasdietrich/todo/model"
)

type database struct {
	*gorm.DB
}

func open(file string) (*database, error) {
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		return nil, errors.Wrap(err, "could not open database")
	}

	return &database{db}, db.AutoMigrate(&model.Task{}).Error
}

func (d *database) pending() *database {
	return &database{DB: d.
		Where("archived is null").
		Order("due asc, id asc"),
	}
}
