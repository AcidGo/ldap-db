import logging
import pymysql
from pyzabbix import ZabbixAPI

class ZBXAPI(object):
    def __init__(self, host):
        self.zapi = ZabbixAPI(host)

        self.user = None

    def login(self, user, passwd):
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
        logging.info("doing update_user with params: {!s}".format(str(params)))
        self.zapi.user.update(**params)

    def create_user(self, alias, passwd, usrgrp_list, name, surname):
        usrgrps = []
        usrgrps_all = self.zapi.usergroup.get(output="extend")
        for i in usrgrps_all:
            if i.get("name", "") in usrgrp_list:
                usrgrps.append({"usrgrpid": i["usrgrpid"]})
        logging.info(f"doing create_user with params: alias: {alias}, passwd: {passwd}, usrgrp_list: {usrgrp_list}")
        self.zapi.user.create(
            alias = alias,
            passwd = passwd,
            usrgrps = usrgrps,
            name = name,
            surname = surname,
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
        with self.db.cursor(cursor=pymysql.cursors.DictCursor) as cur:
            cur.execute(query)
            res = cur.fetchall()
        return res

class ZabbixUserSync(object):
    def __init__(self, zbx_api):
        self.zapi = zbx_api

        if self.zapi is None:
            raise Exception("the self items is nil")

    def set_default_usrgrp_list(self, usrgrp_list):
        self.usrgrp_list = usrgrp_list

    def compare(self, db_users, zbx_users):
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
                self.zapi.create_user(i["alias"], i.get("passwd", ""), self.usrgrp_list, i.get("name", ""), i.get("surname", ""))

def init_logger(level, logfile=None):
    """日志功能初始化。
    如果使用日志文件记录，那么则默认使用 RotatingFileHandler 的大小轮询方式，
    默认每个最大 10 MB，最多保留 5 个。
    Args:
        level: 设定的最低日志级别。
        logfile: 设置日志文件路径，如果不设置则表示将日志输出于标准输出。
    """
    import os
    import sys
    from logging.handlers import RotatingFileHandler
    if not logfile:
        logging.basicConfig(
            level = getattr(logging, level.upper()),
            format = "%(asctime)s [%(levelname)s] %(message)s",
            datefmt = "%Y-%m-%d %H:%M:%S"
        )
    else:
        logger = logging.getLogger()
        logger.setLevel(getattr(logging, level.upper()))
        if logfile.lower() == "local":
            logfile = os.path.join(sys.path[0], os.path.basename(os.path.splitext(__file__)[0]) + ".log")
        handler = RotatingFileHandler(logfile, maxBytes=10*1024*1024, backupCount=5)
        formatter = logging.Formatter("%(asctime)s [%(levelname)s] %(message)s", "%Y-%m-%d %H:%M:%S")
        handler.setFormatter(formatter)
        logger.addHandler(handler)
    logging.info("Logger init finished.")

if __name__ == "__main__":
    from config import *

    init_logger("info")

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
    zbx_user_sync.set_default_usrgrp_list(["Zabbix administrators"])

    db_users = db_conn.select_dict(db_sql)
    zbx_users = zbx_api.get_all_user()
    zbx_user_sync.compare(db_users, zbx_users)