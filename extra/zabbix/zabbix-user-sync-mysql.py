import pymysql
from pyzabbix import ZabbixAPI

class ZBXAPI(object):
    def __init__(self, host):
        self.zapi = Zabbix(host)

        self.user = None

    def login(self, user, passwd)
        self.user = user
        self.zapi.login(user, passwd)

    def get_all_user(self):
        res = {}
        for i in self.zapi.user.get(output="extend"):
            res[i["alias"]] = {
                "userid": i["userid"],
                "name": i["name"],
                "surname": i["surname"],
            }
        return res

    def update_user(self, alias, name, surname):
        userid = None
        for i in self.zapi.user.get(userid=userid, output="extend"):
            if i["alias"] == alias:
                userid = i["userid"]
                break

        params = {"userid": userid}
        if alias:
            params["alias"] = alias
        if name:
            params["name"] = name
        if surname:
            params["surname"] = surname
        self.zapi.user.update(**params)

    def create_user(self, alias, passwd, usrgrp_list):
        usrgrps = []
        usrgrps_all = self.zapi.usergroup.get(output="extend")
        for i in usrgrps_all:
            if i.get("name", "") in usrgrp_list:
                usrgrps.append({"usrgrpid": i["usrgrpid"]})

        self.zapi.user.create(
            alias = alias,
            passwd = passwd,
            usrgrps = usrgrps,
        )

class MySQLDB(object):
    def __init__(self, host, port, user, passwd, db, charset="utf8"):
        self.db = pymysql.connect(
            host = host,
            port = port,
            user = user,
            password = passwd,
            database = db,
            charset = charset,
        )
        self.db.ping()

    def select_dict(self, query):
        res = {}
        with self.db.cursor() as cur:
            cur.execute(query)
            res = cursor.fetchall()
        return res

class ZabbixUserSync(object):
    def __init__(self, zbx_api):
        self.zapi = zbx_api

        if self.zapi is None:
            raise Exception("the self items is nil")

    def compare(db_users, zbx_users):
        for i in db_users:
            if i["alias"] in zbx_users:
                zi = zbx_users[i["alias"]]
                need_cover = False
                if "name" in i and i["name"] and i["name"] != zi["name"]:
                    need_cover = True
                if "surname" in i and i["surname"] and i["surname"] != zi["surname"]:
                    need_cover = True
                if not need_cover:
                    continue
                self.zapi.update_user(i["alias"], i.get("name", ""), i.get("surname", ""))
            else:
                self.zapi.create_user(i["alias"], i.get("passwd", ""), i.get("surname", ""))

if __name__ == "__main__":
    # TODO(20210401-AcidGo): get config from file
    db_host = ""
    db_port = 3306
    db_user = ""
    db_passwd = ""
    db_db = ""
    db_charset = ""
    db_sql = ""
    zbx_url = ""
    zbx_user = ""
    zbx_passwd = ""

    zbx_api = ZBXAPI(zbx_url)
    zbx_api.login(zbx_user, zbx_passwd)

    db_conn = MySQLDB(
        host = db_host,
        port = db_port,
        user = db_user,
        passwd = db_passwd,
        db = db_db,
        charset = db_charset,
    )

    zbx_user_sync = ZabbixUserSync(zbx_api)

    db_users = db_conn.select_dict(db_sql)
    zbx_users = zbx_api.get_all_user()
    zbx_user_sync.compare(db_users, zbx_users)