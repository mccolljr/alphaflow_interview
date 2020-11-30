package models

import (
	"database/sql"

	"github.com/attache/attache"
)

type Subscription struct {
	ID     int64  `db:"id"`
	UserID int64  `db:"user_id"`
	Pair   string `db:"pair"`
}

func NewSubscription() attache.Record { return new(Subscription) }

func (m *Subscription) Table() string { return "subscriptions" }

func (m *Subscription) Key() (columns []string, values []interface{}) {
	columns = []string{"id"}
	values = []interface{}{m.ID}
	return
}

func (m *Subscription) Insert() (columns []string, values []interface{}) {
	columns = []string{"user_id", "pair"}
	values = []interface{}{m.UserID, m.Pair}
	return
}

func (m *Subscription) AfterInsert(result sql.Result) {
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	m.ID = id
}

func (m *Subscription) Update() (columns []string, values []interface{}) {
	columns = []string{"user_id", "pair"}
	values = []interface{}{m.UserID, m.Pair}
	return
}
