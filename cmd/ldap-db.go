package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "github.com/AcidGo/ldap-db/db"
    "github.com/AcidGo/ldap-db/logger"
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
    SECTION_MAIN_KEY_LOG_DIR        = "log_dir"
    SECTION_MAIN_KEY_LOG_NAME       = "log_name"
    SECTION_MAIN_KEY_LOG_LEVEL      = "log_level"
    SECTION_MAIN_KEY_LOG_REPORT     = "log_report"
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
    logDir              string
    logName             string
    logLevel            string
    logReport           bool

    // runtime
    cfgPath             string
    ldb                 *db.DBConn
    lsvr                *server.Server

    // logger
    logging             *logger.ContextLogger

    // app info
    AppName             string
    AppAuthor           string
    AppVersion          string
    AppGitCommitHash    string
    AppBuildTime        string
    AppGoVersion        string
)

func init() {
    // init logger
    logging = logger.FitContext("ldap-db")

    // parse flag
    flag.StringVar(&cfgPath, "f", "ldap-db.ini", "config file path")
    flag.Usage = flagUsage
    flag.Parse()

    cfg, err := ini.Load(cfgPath)
    if err != nil {
        logging.Fatal(err)
    }

    sec, err := cfg.GetSection(SECTION_MAIN)
    if err != nil {
        logging.Fatal(err)
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
    logDir          = sec.Key(SECTION_MAIN_KEY_LOG_DIR).String()
    logName         = sec.Key(SECTION_MAIN_KEY_LOG_NAME).String()
    logLevel        = sec.Key(SECTION_MAIN_KEY_LOG_LEVEL).String()
    logReport, err  = sec.Key(SECTION_MAIN_KEY_LOG_REPORT).Bool()

    if !IsDir(logDir) {
        logging.Fatalf("the log dir %s is not a dir or no exists", logDir)
    }
    if err != nil {
        logging.Fatal(err)
    }

    logPath := filepath.Join(logDir, logName)
    err = logger.LogFileSetting(logPath)
    if err != nil {
        logging.Fatal(err)
    }
    logger.ReportCallerSetting(logReport)
    err = logger.LogLevelSetting(logLevel)
    if err != nil {
        logging.Fatal(err)
    }
}

func main() {
    var err error
    // create db conn to backend source
    ldb, err = db.NewDBConn(dbDriver, dbDsn)
    if err != nil {
        logging.Fatal(err)
    }

    // create ldap server
    lsvr, err = server.NewServer(ldb, serverListen)
    if err != nil {
        logging.Fatal(err)
    }
    // setting for ldap server
    err = lsvr.SetBind(ldapBindDn, ldapBindPasswd)
    if err != nil {
        logging.Fatal(err)
    }
    err = lsvr.SetBase(ldapBaseDn, ldapBaseEn, ldapBaseQuery, ldapBaseCrypt)
    if err != nil {
        logging.Fatal(err)
    }

    // running ldap server
    err = lsvr.ListenAndServe()
    if err != nil {
        logging.Fatal(err)
    }
}

func IsDir(path string) (bool) {
    s, err := os.Stat(path)
    if err != nil {
        return false
    }
    return s.IsDir()
}

func flagUsage() {
    usageMsg := fmt.Sprintf(`App: %s
Version: %s
Author: %s
GitCommit: %s
BuildTime: %s
GoVersion: %s
Options:
`, AppName, AppVersion, AppAuthor, AppGitCommitHash, AppBuildTime, AppGoVersion)

    fmt.Fprintf(os.Stderr, usageMsg)
    flag.PrintDefaults()
}