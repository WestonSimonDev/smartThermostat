package Helpers

import (
	"fmt"
	"math"
	"thermostat/Types"
	"time"
)

func CalculateThermostatActions(currentState Types.ThermostatState, lastCommand Types.StoredThermostatAction) Types.ThermostatAction {
	tempDiff := int(math.Round(float64(currentState.IndicatedTemp))) - currentState.SetTemp
	var action Types.ThermostatAction
	if (currentState.SetTemp == int(currentState.IndicatedTemp)) && ((currentState.Heat == false) && (currentState.Cool == false) && (currentState.Blower == false)) {
		action.Execute = false
		action.Heat = false
		action.Cool = false
		action.Blower = false
	} else {

		if math.Abs(float64(tempDiff)) >= 2 {

			if tempDiff < 0 {
				//turn on heater if i-temp is less than set temp
				action.Execute = true
				action.Cool = false
				action.Blower = true
				action.Heat = true
			} else if tempDiff > 0 {
				//turn on ac if i-temp is more than set temp
				action.Execute = true
				action.Cool = true
				action.Blower = true
				action.Heat = false
			}

		} else if tempDiff == 0 {
			action.Execute = (time.Now().Sub(*lastCommand.Timestamp) >= 5*time.Minute)
			action.Cool = false
			action.Heat = false
			action.Blower = false

		} else if (math.Abs(float64(tempDiff)) < 2) && (time.Now().Sub(*lastCommand.Timestamp) >= 5*time.Minute) {
			if tempDiff < 0 {
				action.Execute = true
				action.Cool = false
				action.Blower = true
				action.Heat = true
			} else if tempDiff > 0 {
				//turn on ac if i-temp is more than set temp
				action.Execute = true
				action.Cool = true
				action.Blower = true
				action.Heat = false
			}
		} else if (math.Abs(float64(tempDiff)) < 2) && !(time.Now().Sub(*lastCommand.Timestamp) >= 5*time.Minute) {
			action.Execute = true
			action.Cool = currentState.Cool
			action.Heat = currentState.Heat
			action.Blower = currentState.Blower
		}
	}

	if action.Execute == true {
		if action.Heat == currentState.Heat && action.Cool == currentState.Cool && action.Blower == currentState.Blower {
			action.Execute = false
		}
	}
	fmt.Println((time.Now().Sub(*lastCommand.Timestamp) >= 5*time.Minute))
	fmt.Println(time.Now().Sub(*lastCommand.Timestamp))
	return action

}
