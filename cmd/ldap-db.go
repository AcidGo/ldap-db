package main

const (
    SECTION_MAIN                = "main"
    SECTION_MAIN_KEY_LISTEN     = "listen"
    SECTION_MAIN_KEY_BINDUSER   = "bind_user"
    SECTION_MAIN_KEY_BINDPASSWD = "bind_passwd"
    SECTION_MAIN_KEY_QUERYSTMT  = "query_stmt"
)

var (
    // config
    serverListen        string
    ldapBindUser        string
    ldapBindPasswd      string
    ldapQueryStmt       string

    // runtime
    cfgPath             string
    db                  *db.DBConn
)

func init() {
    // parse flag
    flag.StringVar(&cfgPath, "f", "ldap-db", "config file path")
    flag.Parse()

    cfg, err := ini.Load(cfgPath)
    if err != nil {
        log.Fatal(err)
    }

    sec, err := cfg.GetSection(SECTION_MAIN)
    if err != nil {
        log.Fatal(err)
    }

    serverListen    := sec.Key(SECTION_MAIN_KEY_LISTEN).String()
    ldapBindUser    := sec.Key(SECTION_MAIN_KEY_BINDUSER).String()
    ldapBindPasswd  := sec.Key(SECTION_MAIN_KEY_BINDPASSWD).String()
    ldapQueryStmt   := sec.Key(SECTION_MAIN_KEY_QUERYSTMT).String()
}

func main() {
    // create a new LDAP server
    
}

