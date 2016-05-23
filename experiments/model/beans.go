package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int       `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	FirstName   string    `sql:"size:255;index:name_idx"`
	LastName    string    `sql:"size:255;index:name_idx"`
	DateOfBirth time.Time `sql:"DEFAULT:current_timestamp"`
	Address     *Address  // one-to-one relationship
	Deleted     bool      `sql:"DEFAULT:false"`
}

type Address struct {
	street   sql.NullString
	area     sql.NullString
	pin_code string `gorm:"primary_key"`
	City     string `gorm:"primary_key"`
	Country  string `gorm:"primary_key"`
}

type Company struct {
	ID      int     `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name    string  `sql:"size:255;unique;index"`
	Users   []User  // one-to-many relationship
	Address Address // one-to-one relationship
}
/*
type New struct {
	gorm.Model
	Code string
}

type Productss struct {
	gorm.Model
}*/

type Person struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Deleted   bool   `db:"deleted"`
	CityId    int    `db:"city_id"`
	City      *City
}

type City struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Employee struct {
	Name string
	Id   int64
	// BossId is an id into the employee table
	Boss_Id sql.NullInt64 `db:"boss_id"`
}

/*type Person struct {
	ID        int    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	FirstName string `sql:"size:255;index:name_idx"`
	LastName  string `sql:"size:255;index:name_idx"`
	Deleted   bool   `sql:"DEFAULT:false"`
	City      *City
}
	model.Person{}

type City struct {
	ID   int    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name string `sql:"size:255;unique;index"`
}*/
