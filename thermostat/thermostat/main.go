package main

import (
	"fmt"
	"thermostat/Connector"
	"thermostat/Helpers"
	"time"
)

func main() {
	Connector.InitDB()

	for {
		currentState := Helpers.GetThermostatState()
		lastCommand := Helpers.GetLastCommand()
		action := Helpers.CalculateThermostatActions(currentState, lastCommand)
		if action.Execute {
			executorResponse := Helpers.CommandExecutor(action)
			if executorResponse != "successful" {
				fmt.Println("Executor error: ", executorResponse)
			}
		}
		fmt.Printf("C state: %+v \n", currentState)
		fmt.Printf("L state: %+v \n", lastCommand)
		fmt.Printf("Action: %+v \n", action)
		time.Sleep(500 * time.Millisecond)

	}
}
