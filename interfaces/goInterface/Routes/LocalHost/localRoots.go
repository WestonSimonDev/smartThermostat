package LocalHost

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"goInterface/Connector"
	"goInterface/Routes/LocalHost/Types"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type requestBody struct {
	Action   string      `json:"action"`
	Payload  interface{} `json:"payload"`
	Response *string     `json:"response"`
}

type heartBeatPayload struct {
	Change   bool                   `json:"change"`
	NewState *Types.ThermostatState `json:"newState"`
}

type changeTempsPayload struct {
	Direction int `json:"direction"`
}

type changeModePayload struct {
	Mode string `json:"mode"`
}

type changeFanPayload struct {
	Fan string `json:"fan"`
}

type basicErrorReponse struct {
	Error string `json:"error"`
}

func getThermostatState() (Types.ThermostatState, error) {
	var state Types.ThermostatState
	indicatedTemp := 0
	state.IndicatedTemp = &indicatedTemp
	err := Connector.DB.QueryRow("SELECT setTemp, CAST(indicatedTemp AS SIGNED INTEGER), heat, cooling, blower, mode, fan FROM thermostatProperties").Scan(&state.SetTemp, &state.IndicatedTemp, &state.Heat, &state.Cool, &state.Blower, &state.Mode, &state.Fan)
	if err != nil {
		return Types.ThermostatState{}, err
	} else {
		return state, nil
	}

}

func checkIfStateMatches(s1 Types.ThermostatState, s2 Types.ThermostatState) bool {
	setTemp := (s1.SetTemp == s2.SetTemp)
	iTemp := (s1.IndicatedTemp != nil && s2.IndicatedTemp != nil && *s1.IndicatedTemp == *s2.IndicatedTemp)
	heat := (s1.Heat == s2.Heat)
	cool := (s1.Cool == s2.Cool)
	blower := (s1.Blower == s2.Blower)
	mode := (s1.Mode == s2.Mode)
	fan := (s1.Fan == s2.Fan)

	return setTemp && iTemp && heat && cool && blower && mode && fan
}

func requestRouter(request requestBody, lastState Types.ThermostatState) (requestBody, Types.ThermostatState) {
	if request.Action == "heartBeat" {
		currentState, stateErr := getThermostatState()
		if stateErr != nil {
			fmt.Println(stateErr)
		} else {
			if !checkIfStateMatches(currentState, lastState) {

				lastState = currentState
				return requestBody{Action: "response/heartBeat", Payload: heartBeatPayload{Change: true, NewState: &lastState}}, lastState
			} else {

				return requestBody{Action: "response/heartBeat", Payload: heartBeatPayload{Change: false}}, lastState
			}
		}
		return requestBody{Action: "response/heartBeat"}, lastState
	}

	var slices = strings.Split(request.Action, "/")
	fmt.Println(slices)
	switch slices[0] {
	case "temps":
		switch slices[1] {
		case "change":
			fmt.Println("hi")
			fmt.Printf("%+v \n", request)
			payloadMap, mapError := request.Payload.(map[string]interface{})
			if !mapError {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}
			var payload changeTempsPayload
			payloadBytes, jsonErr := json.Marshal(payloadMap)
			if jsonErr != nil {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			jsonMapErr := json.Unmarshal(payloadBytes, &payload)
			if jsonMapErr != nil {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			_, dbErr := Connector.DB.Exec("UPDATE thermostatProperties SET setTemp = setTemp + ?", payload.Direction)
			if dbErr != nil {
				fmt.Println(dbErr)
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "DB ERROR"}}, Types.ThermostatState{}
			} else {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "successful"}}, Types.ThermostatState{}
			}

		}
	case "mode":
		switch slices[1] {
		case "change":
			fmt.Println("hi")
			fmt.Printf("%+v \n", request)
			payloadMap, mapError := request.Payload.(map[string]interface{})

			if !mapError {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			var payload changeModePayload
			payloadBytes, jsonErr := json.Marshal(payloadMap)

			if jsonErr != nil {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			jsonMapErr := json.Unmarshal(payloadBytes, &payload)

			if jsonMapErr != nil {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			_, dbErr := Connector.DB.Exec("UPDATE thermostatProperties SET mode = ?", payload.Mode)

			if dbErr != nil {
				fmt.Println(dbErr)
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "DB ERROR"}}, Types.ThermostatState{}
			} else {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "successful"}}, Types.ThermostatState{}
			}

		}
	case "fan":
		switch slices[1] {
		case "change":
			fmt.Println("hi")
			fmt.Printf("%+v \n", request)
			payloadMap, mapError := request.Payload.(map[string]interface{})

			if !mapError {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			var payload changeFanPayload
			payloadBytes, jsonErr := json.Marshal(payloadMap)

			if jsonErr != nil {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			jsonMapErr := json.Unmarshal(payloadBytes, &payload)

			if jsonMapErr != nil {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "Invalid Payload"}}, Types.ThermostatState{}
			}

			_, dbErr := Connector.DB.Exec("UPDATE thermostatProperties SET fan = ?", payload.Fan)

			if dbErr != nil {
				fmt.Println(dbErr)
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "DB ERROR"}}, Types.ThermostatState{}
			} else {
				return requestBody{Action: fmt.Sprintf("response/%s", *request.Response), Payload: basicErrorReponse{Error: "successful"}}, Types.ThermostatState{}
			}

		}
	}

	return requestBody{Action: "response/some"}, lastState
}

func InitLocalRouts(router *mux.Router) {
	router.HandleFunc("/ws", func(output http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(output, req, nil)
		if err != nil {
			fmt.Println("Some error: ", err)
			return
		}

		var IndicatedTemp = 0
		var lastState = Types.ThermostatState{IndicatedTemp: &IndicatedTemp, SetTemp: 0, Heat: false, Cool: false, Blower: false}

		for {
			msgType, msgBody, msgErr := conn.ReadMessage()
			if msgErr != nil {
				fmt.Println("Error reading message", msgErr)
				conn.Close()
				return
			}
			if msgType != 1 {
				fmt.Println("Error bad msg type")
				conn.Close()
				return
			}

			var requestData requestBody

			jsonErr := json.Unmarshal([]byte(msgBody), &requestData)
			if jsonErr != nil {
				fmt.Println("Json error: ", jsonErr)
				msgString := string(msgBody)
				conn.WriteJSON(requestBody{Action: "error/json", Response: &msgString})
			} else {
				//fmt.Println(string(msgBody))
				//fmt.Printf("%+v \n", requestData)
				responseMsg, updatedState := requestRouter(requestData, lastState)
				lastState = updatedState
				conn.WriteJSON(responseMsg)
			}

		}

	})
}
