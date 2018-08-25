package request

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type Validater interface {
	Validate() error
}

func DecodeAndValidateJSON(bodyIn io.Reader, dest Validater) error {
	if err := json.NewDecoder(bodyIn).Decode(dest); err != nil {
		return errors.Wrapf(err, "Failed to decode body as json")
	}

	if err := dest.Validate(); err != nil {
		return err
	}

	return nil
}
