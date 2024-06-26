package Helpers

import (
	"fmt"
	"thermostat/Connector"
	"thermostat/Types"
)

func GetThermostatState() Types.ThermostatState {
	var state Types.ThermostatState
	err := Connector.DB.QueryRow("SELECT setTemp, indicatedTemp, heat, cooling, blower, mode, fan FROM thermostatProperties LIMIT 1;").Scan(&state.SetTemp, &state.IndicatedTemp, &state.Heat, &state.Cool, &state.Blower, &state.Mode, &state.Fan)
	if err != nil {
		fmt.Println("Get state error:", err)
	}
	return state
}
