package models

import (
	"database/sql"

	"github.com/attache/attache"
)

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func NewUser() attache.Record { return new(User) }

func (m *User) Table() string { return "users" }

func (m *User) Key() (columns []string, values []interface{}) {
	columns = []string{"id"}
	values = []interface{}{m.ID}
	return
}

func (m *User) Insert() (columns []string, values []interface{}) {
	columns = []string{"name"}
	values = []interface{}{m.Name}
	return
}

func (m *User) AfterInsert(result sql.Result) {
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	m.ID = id
}

func (m *User) Update() (columns []string, values []interface{}) {
	columns = []string{"name"}
	values = []interface{}{m.Name}
	return
}
