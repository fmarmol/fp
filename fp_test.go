package fp

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func TestParseMyInt(t *testing.T) {
	var i MyInt
	err := parseString(&i, "-1")
	require.NoError(t, err)
	require.Equal(t, MyInt(0), i)
}

type MyCustomType struct{}

func TestParseMyCustomeType(t *testing.T) {
	var i MyCustomType
	err := parseString(&i, "-1")
	require.Error(t, err)
}

func TestParseUint16(t *testing.T) {
	var i uint16
	err := parseString(&i, "1")
	require.NoError(t, err)
	require.Equal(t, uint16(1), i)
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalText(text []byte) error {
	t, err := time.Parse(time.DateOnly, string(text))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func TestParse(t *testing.T) {

	type SubContent struct {
		Floats []float32 `fp:"floats"`
	}

	type Content struct {
		String     string   `fp:"string"`
		Strings    []string `fp:"strings"`
		Ints       []int    `fp:"ints"`
		Dates      []Date   `fp:"dates"`
		SubContent SubContent
		Default    uint8   `fp:"default" fp-def:"42"`
		Defaults   []int16 `fp:"defaults" fp-def:"-1"`
	}

	values := url.Values{
		"string":  []string{"bob"},
		"strings": []string{"a", "b"},
		"ints":    []string{"1", "2", "3"},
		"floats":  []string{"1.1", "2.2", "3.3"},
		"dates":   []string{"2024-01-01", "2024-02-02", "2024-03-03"},
	}
	var content Content
	err := Parse(&content, values)
	require.NoError(t, err)

	require.Equal(t, values["string"][0], content.String)
	require.Equal(t, values["strings"], content.Strings)
	require.Equal(t, []int{1, 2, 3}, content.Ints)

	for index, date := range content.Dates {
		require.Equal(t, values["dates"][index], date.Format(time.DateOnly))
	}
	require.Equal(t, []float32{1.1, 2.2, 3.3}, content.SubContent.Floats)
	require.Equal(t, uint8(42), content.Default)
	require.Equal(t, []int16{-1}, content.Defaults)
}
