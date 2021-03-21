package server

import (
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/AcidGo/ldap-db/db"
    "gopkg.in/ini.v1"
    ldap "github.com/vjeantet/ldapserver"
)

const (
    QUERY_BASE_DN   = "dn=ldap-db"
    QUERY_LABEL     = "Search - LDAP DB"
)

type Server struct {
    lSvr        *ldap.Server
    dDB         *db.DBConn
    listen      string
    bindUser    string
    bindPasswd  string
}

func NewServer(db *db.DBConn, l, bUser, bPasswd string) (*Server, error) {
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
        bindUser:   bUser,
        bindPasswd: bPasswd,
    }

    return server, nil
}

func (svr *Server) Handle() (error) {
    routes := ldap.NewRouteMux()
    routes.NotFound(svr.handleNotFound)
    routes.Bind(svr.handleBind)
    rotues.Search(svr.handleSearchLdapDb)

    svr.lSvr.Handle(routes)

    return nil
}

func (svr *Server) ListenAndServe() (error) {
    return svr.lSvr.ListenAndServe(svr.listen)
}

func (svr *Server) handleBind(w ldap.ResponseWriter, m *ldap.Message) {
    r := m.GetBindRequest()
    res := ldap.NewBindResponse(ldap.LDAPResultSuccess)

    if r.AuthenticationChoice() == "simple" {
        if string(r.Name()) == svr.bindUser && string(r.Authentication()) == svr.bindPasswd {
            w.Write(res)
            return
        }

        log.Printf("Bind failed User=%s, Pass=%#v\n", string(r.Name()), r.Authentication())
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

func (svr *Server) handleSearchLdapDb(w ldap.ResponseWriter, r *ldap.Message) {
    r := m.GetSearchRequest()
    log.Printf("Request BaseDn=%s\n", r.BaseObject())
    log.Printf("Request Filter=%s\n", r.Filter())
    log.Printf("Request FilterString=%s\n", r.FilterString())
    log.Printf("Request Attributes=%s\n", r.Attributes())
    log.Printf("Request TimeLimit=%d\n", r.TimeLimit().Int())

    // Handle Stop Signal (server stop / client disconnected / Abandoned request....)
    select {
    case <-m.Done:
        log.Print("Leaving handleSearch...")
        return
    default:
    }

    e := ldap.NewSearchResultEntry(fmt.Sprintf("%s, %s", QUERY_BASE_DN, string(r.BaseObject())))
    e.AddAttribute("sn", "0612324567")
    e.AddAttribute("telephoneNumber", "0612324567")
    e.AddAttribute("cn", "Valère JEANTET")
    w.Write(e)

    e = ldap.NewSearchResultEntry("cn=Claire Thomas, " + string(r.BaseObject()))
    e.AddAttribute("mail", "claire.thomas@gmail.com")
    e.AddAttribute("cn", "Claire THOMAS")
    w.Write(e)

    res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
    w.Write(res)
}