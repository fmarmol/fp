package fp

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Date2 time.Time

func (d *Date2) UnmarshalText(text []byte) error {
	t, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return err
	}
	*d = Date2(t)
	return nil
}

func TestDate2(t *testing.T) {
	data := map[string][]string{
		"start": {"2024-03-24"},
	}

	type Request struct {
		Start Date2 `fp:"start"`
	}

	var r Request
	err := Parse(&r, data)
	require.NoError(t, err)
	fmt.Println("START:", time.Time(r.Start).Format(time.DateOnly))

}
