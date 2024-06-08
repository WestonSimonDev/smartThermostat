/*
Author: Weston Simon
Email: weston@wcloud.com

Creation Date: 2024-04-28 12:51:59

© wcloud
*/
package Connector

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goInterface/Types"
	"io/ioutil"
	"log"
	"os"
)

var DB *sql.DB // Declare a global variable to hold the database connection

func InitDB() {

	jsonFile, fileErr := os.Open("../../conf.json")

	if fileErr != nil {
		fmt.Println("File Error: ", fileErr)
		return
	}

	defer jsonFile.Close()

	jsonBytes, _ := ioutil.ReadAll(jsonFile)

	var conf Types.ConfStructure

	jsonErr := json.Unmarshal(jsonBytes, &conf)
	if jsonErr != nil {
		fmt.Println("JSON Error: ", jsonErr)
		return
	}

	var err error

	// Open a database connection
	DB, err = sql.Open("mysql", fmt.Sprintf("superSmartThermostat:%s@(192.168.88.160:3306)/temps?parseTime=true", conf.DBPassword))
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	} else {
		fmt.Println("DB Init Complete")
	}

	// Ping the database to ensure connectivity
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	} else {
		fmt.Println("DB Ping Init Complete")
	}

}
