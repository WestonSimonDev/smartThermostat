package Helpers

import (
	"fmt"
	"thermostat/Connector"
	"thermostat/Types"
)

func GetLastCommand() Types.StoredThermostatAction {
	var command Types.StoredThermostatAction
	err := Connector.DB.QueryRow("SELECT timeStamp, heat, cooling, blower FROM thermostatCommands ORDER BY pID DESC LIMIT 1;").Scan(&command.Timestamp, &command.Heat, &command.Cool, &command.Blower)
	if err != nil {
		fmt.Println("Get command error:", err)
	}
	return command
}
