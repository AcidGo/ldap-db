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
    def __init__(self, host, port, user, passwd, db, charset):
        self.db = pymysql.connect(

        )

    def select_dict(self, query):
        res = {}
        with self.db.cursor() as cur:
            cur.execute(query)
            res = cursor.fetchall()
        return res

class ZabbixUserSyncMySQL(object):
    def

    def __init__(self, zbx_api, mysql_db):
        self.zapi = zbx_api
        self.mydb = mysql_db

        if self.zapi is None or self.mydb is None:
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
                self.zapi.create_user(i["alias"], i.get("passwd", ""), )