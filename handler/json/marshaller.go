package json

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// UnmarshalJSON unmarshals the given string into the given struct
func UnmarshalJSON(s string, t interface{}) error {
	if len(s) == 0 {
		return errors.New("Cannot unmarshal empty string")
	}

	if t == nil {
		return errors.New("UnmarshalJSON needs a non-nil interface to unmarshal into")
	}

	err := json.Unmarshal([]byte(s), &t)
	if err != nil {
		return errors.Wrapf(err, "Could not unmarshal string %v", s)
	}

	return nil
}
