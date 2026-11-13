package models

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type Runtime time.Duration

var RuntimeMarshaler = json.MarshalToFunc(
	func(enc *jsontext.Encoder, val Runtime) error {
		mins := math.Round(time.Duration(val).Minutes())
		token := jsontext.String(fmt.Sprintf("%dm", int(mins)))
		return enc.WriteToken(token)
	},
)

var RuntimeUnmarshaler = json.UnmarshalFromFunc(
	func(dec *jsontext.Decoder, val *Runtime) error {
		token, err := dec.ReadToken()
		if err != nil {
			return err
		}

		if token.Kind() != jsontext.KindString || !strings.HasSuffix(token.String(), "m") {
			return errors.New("invalid runtime value")
		}

		d, err := time.ParseDuration(token.String())
		if err != nil {
			return fmt.Errorf("invalid runtime value: %w", err)
		}

		*val = Runtime(d)
		return nil
	},
)
