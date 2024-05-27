import mysql.connector
from mysql.connector import Error
import mysql.connector.pooling
import time
import json

with open("../../conf.json", "r") as outfile:
    password = json.load(outfile)["dbPassword"]

print("Connector started")
dbconfig = {"host":'localhost',
            "database":'temps',
            "user":'superSmartThermostat',
            "password": password,
            "auth_plugin":'mysql_native_password'
            }

cnx = mysql.connector.pooling.MySQLConnectionPool(pool_name = "mypool",
                                                  pool_size = 30, #18, #scaled up for delivery report. Attempting to prevent self DOS attack. Max concurrent connections 72.
                                                  **dbconfig,
                                                  use_unicode=True,
                                                  charset='utf8mb4',
                                                  collation='utf8mb4_general_ci')

print("Connection established")

conn_count = 0
class conn_controller:


    def get_db_conn():
        global conn_count
        pool_size = cnx.pool_size
        while not pool_size - 3 > conn_count:
            time.sleep(0.000001)
        print("Creating DB connection")
        conn_count += 1
        print(conn_count)
        return cnx.get_connection()

    def close_conn(db):
        global conn_count
        print("Closing DB connetion")
        db.close()
        conn_count -= 1
        print("DB Connections", conn_count)
