package server

import (
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/AcidGo/ldap-db/db"
    ldap "github.com/vjeantet/ldapserver"
)

const (
    SEARCH_DN_ENTRY_ATTR    = "cn"
    SEARCH_DN_ENTRY_VAL     = "acidgo"
    SEARCH_BIND_ATTR        = "userPassword"
    QUERY_LABEL             = "Search - LDAP DB"
)

type Server struct {
    lSvr        *ldap.Server
    dDB         *db.DBConn
    listen      string
    bindDn      string
    bindPasswd  string
    baseDn      string
    baseEn      string
    baseQuery   string
    baseCrypt   string
}

func NewServer(db *db.DBConn, l string) (*Server, error) {
    if db == nil {
        return nil, fmt.Errorf("the backend db is nil")
    }

    if l == "" {
        return nil, fmt.Errof("the listen bind is empty")
    }

    lSvr := ldap.NewServer()
    server := &Server{
        lSvr:       lSvr,
        dDB:        db,
        listen:     l,
    }

    return server, nil
}

func (svr *Server) SetBind(bindDn, bindPasswd string) (error) {
    svr.bindDn = bindDn
    svr.bindPasswd = bindPasswd

    return nil
}

func (svr *Server) SetBase(baseDn, baseEn, baseQuery, baseCrypt string) (error) {
    svr.baseDn = baseDn
    svr.baseEn = baseEn
    svr.baseQuery = baseQuery
    svr.baseCrypt = baseCrypt

    return nil
}

func (svr *Server) handle() (error) {
    routes := ldap.NewRouteMux()
    routes.NotFound(svr.handleNotFound)
    routes.Bind(svr.handleBind)
    rotues.Search(svr.handleSearch)

    svr.lSvr.Handle(routes)

    return nil
}

func (svr *Server) ListenAndServe() (error) {
    err := svr.handle()
    if err != nil {
        return err
    }

    return svr.lSvr.ListenAndServe(svr.listen)
}

func (svr *Server) handleBind(w ldap.ResponseWriter, m *ldap.Message) {
    r := m.GetBindRequest()
    bName := string(r.Name())
    bAuth := fmt.Sprintf("%v", r.Authentication())
    log.Printf("bind name: %s", bName)
    log.Printf("bind auth: %s", bAuth)

    res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
    if r.AuthenticationChoice() == "simple" {
        if string(r.Name()) == svr.bindDn && 

        if string(r.Name()) == "login" || string(r.Name()) == "cn=testing, ou=Users,dc=acidgo,dc=com" {
            w.Write(res)
            return
        }
        log.Printf("Bind failed User=%s, Pass=%#v", string(r.Name()), r.Authentication())
        res.SetResultCode(ldap.LDAPResultInvalidCredentials)
        res.SetDiagnosticMessage("invalid credentials")
    } else {
        res.SetResultCode(ldap.LDAPResultUnwillingToPerform)
        res.SetDiagnosticMessage("Authentication choice not supported")
    }

    w.Write(res)
}

func (svr *Server) handleNotFound(w ldap.ResponseWriter, r *ldap.Message) {
    switch r.ProtocolOpType() {
    case ldap.ApplicationBindRequest:
        res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
        res.SetDiagnosticMessage("Default binding behavior set to return Success")

        w.Write(res)

    default:
        res := ldap.NewResponse(ldap.LDAPResultUnwillingToPerform)
        res.SetDiagnosticMessage("Operation not implemented by server")
        w.Write(res)
    }
}

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
    r := m.GetSearchRequest()
    log.Printf("Request BaseDn=%s", r.BaseObject())
    log.Printf("Request Filter=%s", r.Filter())
    log.Printf("Request FilterString=%s", r.FilterString())
    log.Printf("Request Attributes=%s", r.Attributes())
    log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())

    // Handle Stop Signal (server stop / client disconnected / Abandoned request....)
    select {
    case <-m.Done:
        log.Print("Leaving handleSearch...")
        return
    default:
    }

    var enVal string
    for _, attr := range r.Attributes() {
        if attr.Type_() == svr.baseEn {
            for _, attrVal := range attr.Vals() {
                if attrVal != "" {
                    enVal = attrVal
                    break
                }
            }
        }
        if enVal != "" {
            break
        }
    }

    e := ldap.NewSearchResultEntry(enVal)
    e.AddAttribute(svr.baseEn, "MOCK")
    e.AddAttribute(SEARCH_DN_ENTRY_ATTR, SEARCH_DN_ENTRY_VAL)
    w.Write(e)

    res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
    w.Write(res)
}