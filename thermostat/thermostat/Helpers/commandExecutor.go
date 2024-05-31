package Helpers

import (
	"thermostat/Connector"
	"thermostat/Types"
	"time"
)

func CommandExecutor(action Types.ThermostatAction) string {
	_, err := Connector.DB.Exec("UPDATE thermostatProperties SET timeStamp = ?, heat = ?, cooling = ?, blower = ?;", time.Now().UTC(), action.Heat, action.Cool, action.Blower)
	if err != nil {
		return err.Error()
	}
	_, err2 := Connector.DB.Exec("INSERT INTO thermostatCommands(timeStamp, heat, cooling, blower) VALUES (?, ?, ?, ?);", time.Now().UTC(), action.Heat, action.Cool, action.Blower)
	if err2 != nil {
		return err2.Error()
	}
	return "successful"
}
