package Types

type ThermostatState struct {
	SetTemp       int    `json:"setTemp"`
	IndicatedTemp *int   `json:"indicatedTemp"`
	Heat          bool   `json:"heat"`
	Cool          bool   `json:"cool"`
	Blower        bool   `json:"blower"`
	Mode          string `json:"mode"`
	Fan           string `json:"fan"`
}
