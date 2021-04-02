package main

import (
    "flag"
    "log"

    "github.com/AcidGo/ldap-db/db"
    "github.com/AcidGo/ldap-db/server"
    "gopkg.in/ini.v1"
)

const (
    SECTION_MAIN                    = "main"
    SECTION_MAIN_KEY_LISTEN         = "listen"
    SECTION_MAIN_KEY_BIND_DN        = "bind_dn"
    SECTION_MAIN_KEY_BIND_PASSWD    = "bind_passwd"
    SECTION_MAIN_KEY_BASE_DN        = "base_dn"
    SECTION_MAIN_KEY_BASE_EN        = "base_en"
    SECTION_MAIN_KEY_BASE_QUERY     = "base_query"
    SECTION_MAIN_KEY_BASE_CRYPT     = "base_crypt"
    SECTION_MAIN_KEY_DB_DRIVER      = "db_driver"
    SECTION_MAIN_KEY_DB_DSN         = "db_dsn"
)

var (
    // config
    serverListen        string
    ldapBindDn          string
    ldapBindPasswd      string
    ldapBaseDn          string
    ldapBaseEn          string
    ldapBaseQuery       string
    ldapBaseCrypt       string
    dbDriver            string
    dbDsn               string

    // runtime
    cfgPath             string
    ldb                 *db.DBConn
    lsvr                *server.Server
)

func init() {
    // parse flag
    flag.StringVar(&cfgPath, "f", "ldap-db.ini", "config file path")
    flag.Parse()

    cfg, err := ini.Load(cfgPath)
    if err != nil {
        log.Fatal(err)
    }

    sec, err := cfg.GetSection(SECTION_MAIN)
    if err != nil {
        log.Fatal(err)
    }

    serverListen    = sec.Key(SECTION_MAIN_KEY_LISTEN).String()
    ldapBindDn      = sec.Key(SECTION_MAIN_KEY_BIND_DN).String()
    ldapBindPasswd  = sec.Key(SECTION_MAIN_KEY_BIND_PASSWD).String()
    ldapBaseDn      = sec.Key(SECTION_MAIN_KEY_BASE_DN).String()
    ldapBaseEn      = sec.Key(SECTION_MAIN_KEY_BASE_EN).String()
    ldapBaseQuery   = sec.Key(SECTION_MAIN_KEY_BASE_QUERY).String()
    ldapBaseCrypt   = sec.Key(SECTION_MAIN_KEY_BASE_CRYPT).String()
    dbDriver        = sec.Key(SECTION_MAIN_KEY_DB_DRIVER).String()
    dbDsn           = sec.Key(SECTION_MAIN_KEY_DB_DSN).String()
}

func main() {
    var err error
    // create db conn to backend source
    ldb, err = db.NewDBConn(dbDriver, dbDsn)
    if err != nil {
        log.Fatal(err)
    }

    // create ldap server
    lsvr, err = server.NewServer(ldb, serverListen)
    if err != nil {
        log.Fatal(err)
    }
    // setting for ldap server
    err = lsvr.SetBind(ldapBindDn, ldapBindPasswd)
    if err != nil {
        log.Fatal(err)
    }
    err = lsvr.SetBase(ldapBaseDn, ldapBaseEn, ldapBaseQuery, ldapBaseCrypt)
    if err != nil {
        log.Fatal(err)
    }

    // running ldap server
    err = lsvr.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}

