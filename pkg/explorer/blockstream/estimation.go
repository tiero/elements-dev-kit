package blockstream

type estimation struct {
	Two               float64 `json:"2"`
	Three             float64 `json:"3"`
	Four              float64 `json:"4"`
	Six               float64 `json:"6"`
	Ten               float64 `json:"10"`
	Twenty            float64 `json:"20"`
	HundreddFortyFour float64 `json:"144"`
	FiveHundredFour   float64 `json:"504"`
	ThousandEight     float64 `json:"1008"`
}

func (e estimation) Low() float64 {
	return (e.Ten + e.Twenty + e.HundreddFortyFour + e.FiveHundredFour + e.ThousandEight) / 5
}

func (e estimation) Medium() float64 {
	return (e.Three + e.Four + e.Six) / 3
}

func (e estimation) High() float64 {
	return e.Two
}
