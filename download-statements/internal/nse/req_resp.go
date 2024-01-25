package nse

import "codermana.com/go/pkg/asdl/entities"

type Nifty50Data struct {
	Priority int64           `json:"priority"`
	Meta     entities.Script `json:"meta,omitempty"`
}

type Nifty50Resp struct {
	Data []Nifty50Data `json:"data"`
}
