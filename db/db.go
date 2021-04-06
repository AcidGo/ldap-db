package db

import (
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
    "github.com/AcidGo/ldap-db/logger"
)

var logging *logger.ContextLogger

func init() {
    logging = logger.FitContext("database")
}

type DBConn struct {
    driver  string
    db      *sql.DB
}

func NewDBConn(driver, dsn string) (*DBConn, error) {
    db, err := sql.Open(driver, dsn)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    dbConn := &DBConn{
        driver:     driver,
        db:         db,
    }

    return dbConn, err
}

func (dbConn *DBConn) BaseSearch(query, val string) (string, error) {
    row := dbConn.db.QueryRow(query, val)

    var res string
    err := row.Scan(&res)
    if err != nil {
        return res, err
    }

    return res, nil
}