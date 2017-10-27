package redis

import (
	"database/sql"
	"fmt"
)

const (
	identifier = "mysql"
)

type MysqlDoer struct {
	DBName   string
	Username string
	Password string
	Protocol string
}

func (p *MysqlDoer) Identifier() string {
	return identifier
}

func (r *MysqlDoer) Do(addr string) (bool, []string) {
	dsn := fmt.Sprintf("%s%s%s%s%s",
		r.Username,
		":"+r.Password,
		r.Protocol,
		"("+addr+")",
		"/"+r.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return false, nil
	}
	db.Close()
	return true, nil
}
