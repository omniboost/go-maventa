package maventa

import (
	"encoding/json"
	"strconv"
)

type Number float64

func (i *Number) UnmarshalJSON(data []byte) error {
	realNumber := 0.0
	err := json.Unmarshal(data, &realNumber)
	if err == nil {
		*i = Number(realNumber)
		return nil
	}

	// error, so maybe it isn't an int
	str := ""
	err = json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	if str == "" {
		*i = 0
		return nil
	}

	realNumber, err = strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	i2 := Number(realNumber)
	*i = i2
	return nil
}

func (i Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(i))
}
