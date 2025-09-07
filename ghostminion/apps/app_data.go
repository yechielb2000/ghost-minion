package apps

import "encoding/json"

type AppType string
type AppParams any

type AppData struct {
	Id     string          `json:"id"`
	Type   AppType         `json:"type"`
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}

func (ad *AppData) UnmarshalParams(s AppParams) error {
	return json.Unmarshal(ad.Params, s)
}
