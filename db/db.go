package db

import (
    "database/sql"
    "log"
)

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

func (dbConn *DBConn) VerifyUser(query, name string) (bool, error) {
    row, err := dbConn.QueryRow(query, name)
    if err != nil {
        return false, err
    }

    var cnt int
    err = row.Scan(&cnt)
    if err != nil {
        return false, err
    }

    if cnt == 1 {
        return true, nil
    }

    return false, fmt.Errorf("the count of result from verifing user is %d", cnt)
}