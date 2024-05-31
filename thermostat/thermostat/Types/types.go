package Types

import "time"

type ConfStructure struct {
	DBPassword string `json:"dbPassword"`
}

type ThermostatState struct {
	Timestamp     *time.Time
	SetTemp       int
	IndicatedTemp float32
	Heat          bool
	Cool          bool
	Blower        bool
}

type StoredThermostatAction struct {
	Timestamp     *time.Time
	SetTemp       int
	IndicatedTemp float32
	Heat          bool
	Cool          bool
	Blower        bool
}

type ThermostatAction struct {
	Execute bool
	Heat    bool
	Cool    bool
	Blower  bool
}
