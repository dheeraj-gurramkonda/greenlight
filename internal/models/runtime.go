package models

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"math"
	"time"
)

type Runtime time.Duration

var RuntimeMarshaler = json.MarshalToFunc(
	func(enc *jsontext.Encoder, val Runtime) error {
		// Note that we need to convert the val parameter back to its underlying
		// type so that we can call the Minutes() method on it.
		mins := math.Round(time.Duration(val).Minutes())
		token := jsontext.String(fmt.Sprintf("%dm", int(mins)))
		return enc.WriteToken(token)
	},
)
