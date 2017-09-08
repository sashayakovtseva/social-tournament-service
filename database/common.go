package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type closer interface {
	Close()
}

func prepareAndAdd(statements []*sql.Stmt,format string, v ...interface{}) (*sql.Stmt, error) {
	stmt, err := conn.Prepare(fmt.Sprintf(format, v...))
	statements = append(statements, stmt)
	return stmt, err
}

func haveToFailAndClose(closer closer, err error) bool {
	if err != nil {
		closer.Close()
		return true
	}
	return false
}

func checkAndClose(stmt *sql.Stmt) {
	if stmt != nil {
		stmt.Close()
	}
}

func replaceConstraintWithCustom(err, custom error) error {
	if err != nil {
		if err := err.(sqlite3.Error); err.Code == sqlite3.ErrConstraint {
			return custom
		}
		return err
	}
	return nil
}