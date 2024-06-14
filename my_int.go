package fp

import (
	"strconv"
)

// MyInt implements encoding.TextUnmarshaler
type MyInt int32

func (i *MyInt) UnmarshalText(text []byte) error {
	res, err := strconv.ParseInt(string(text), 10, 32)
	if err != nil {
		return err
	}
	*i = MyInt(res + 1) // we increment by 1 to show the method is used
	return nil
}
