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
	"io/ioutil"
	"log"
	"os"
	"thermostat/Types"
)

func createDatabase(dbConn *sql.DB) {
	fmt.Println("Starting DB Structure Alignment")
	dbConn.Exec(`CREATE DATABASE IF NOT EXISTS temps`)
	fmt.Println("Starting Table Alignment")
	_, err := dbConn.Exec(`CREATE TABLE thermostatProperties (
    pid INT(11) NOT NULL AUTO_INCREMENT,
    timeStamp DATETIME NOT NULL DEFAULT current_timestamp(),
    setTemp INT(11) NOT NULL,
    indicatedTemp DECIMAL(5,1) NOT NULL,
    heat TINYINT(1) NOT NULL,
    cooling TINYINT(1) NOT NULL,
    blower TINYINT(1) NOT NULL,
    PRIMARY KEY (pid)
	);
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err2 := dbConn.Exec(`CREATE TABLE thermostatState (
    pID INT(11) NOT NULL AUTO_INCREMENT,
    timeStamp DATETIME NOT NULL DEFAULT current_timestamp(),
    setTemp INT(11) NOT NULL,
    indicatedTemp DECIMAL(5,1) NOT NULL,
    heat TINYINT(1) NOT NULL,
    cooling TINYINT(1) NOT NULL,
    blower TINYINT(1) NOT NULL,
    PRIMARY KEY (pID)
	);
	`)
	if err2 != nil {
		fmt.Println(err)
	}

	_, err3 := dbConn.Exec(`CREATE TABLE thermostatCommands (
    pID INT(11) NOT NULL AUTO_INCREMENT,
    timeStamp DATETIME NOT NULL DEFAULT current_timestamp(),
    heat TINYINT(1) NOT NULL,
    cooling TINYINT(1) NOT NULL,
    blower TINYINT(1) NOT NULL,
    PRIMARY KEY (pID)
	);
	`)
	if err3 != nil {
		fmt.Println(err)
	}

}

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
	DB, err = sql.Open("mysql", fmt.Sprintf("superSmartThermostat:%s@(192.168.88.161:3306)/temps?parseTime=true", conf.DBPassword))
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
	createDatabase(DB)

}
