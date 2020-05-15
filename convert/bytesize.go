package convert

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Bytesize struct {
	Data float64
	Unit string
}

type Unit struct {
	symbol      string
	name        string
	base        uint
	exponential uint
}

var ByteUnits = []Unit{
	// https://en.wikipedia.org/wiki/Byte#Unit_symbol

	// SI + Byte
	{"B", "Byte", 1, 1},
	{"kB", "Kilobyte", 10, 3},
	{"MB", "Megabyte", 10, 6},
	{"GB", "Gigabyte", 10, 9},
	{"TB", "Terabyte", 10, 12},
	{"PB", "Petabyte", 10, 15},
	{"EB", "Exabyte", 10, 18},
	{"ZB", "Zettabyte", 10, 21},
	{"YB", "Yottabyte", 10, 24},

	// IEC
	{"KiB", "Kibibyte", 2, 10},
	{"MiB", "Mebibyte", 2, 20},
	{"GiB", "Gibibyte", 2, 30},
	{"TiB", "Tebibyte", 2, 40},
	{"PiB", "Pebibyte", 2, 50},
	{"EiB", "Exbibyte", 2, 60},
	{"ZiB", "Zebibyte", 2, 70},
	{"YiB", "Yobibyte", 2, 80},
}

var ByteUnitMap map[string]Unit

func init() {
	// build ByteUnitMap
	ByteUnitMap = map[string]Unit{}
	for _, unit := range ByteUnits {
		ByteUnitMap[unit.symbol] = unit
	}
}

func (b *Bytesize) String() string {
	return b.ToHumanReadable()
}

func (b *Bytesize) Dump() {
	fmt.Println("Data: ", b.Data)
	fmt.Println("Unit: ", b.Unit)
}

func (b *Bytesize) cleanUnits(input string) string {
	// cleanup based on UnitMap
	for key, value := range ByteUnitMap {
		if strings.EqualFold(key, input) || strings.EqualFold(value.symbol, input) || strings.EqualFold(value.name, input) {
			input = value.symbol
		}
	}

	return input
}

func ParseBytes(data interface{}) (error, *Bytesize) {
	b := &Bytesize{}

	// given data is int; set it directly
	switch d := data.(type) {
	case string:
		i, err := strconv.ParseFloat(d, 64)
		if err == nil {
			b.Data = i
			b.Unit = "B"
		} else {
			re := regexp.MustCompile(`^\s*(-)?\s*(\d+(?:\.\d+)?)\s*(\w+)\s*$`)
			found := re.FindStringSubmatch(d)
			if len(found) == 0 {
				return fmt.Errorf("could not parse input into number and optional unit: %s", d), nil
			}

			i, err = strconv.ParseFloat(found[1]+found[2], 64)
			if err != nil {
				return fmt.Errorf("could not convert input to float64: %s", d), nil
			}

			b.Data = i
			b.Unit = b.cleanUnits(found[3])
		}
	case float64:
		b.Data = d
		b.Unit = "B"
	case uint64:
		b.Data = float64(d)
		b.Unit = "B"
	case uint:
		b.Data = float64(d)
		b.Unit = "B"
	case int64:
		b.Data = float64(d)
		b.Unit = "B"
	case int:
		b.Data = float64(d)
		b.Unit = "B"
	default:
		return fmt.Errorf("can not handle type %T", d), nil
	}

	return nil, b
}

func (b *Bytesize) calc(targetUnit string) float64 {
	// clean units
	b.Unit = b.cleanUnits(b.Unit)

	// convert given values to bytes
	// example: 1000 MB -> Bytes
	// c := 1000 * match.Pow(10, 6)
	// Result: 1.000.000.000
	c := b.Data * math.Pow(
		float64(ByteUnitMap[b.Unit].base),
		float64(ByteUnitMap[b.Unit].exponential),
	)

	// calculate from bytes to target unit
	// example: Bytes -> Gigabytes
	// x := 1.000.000.000 / (match.Pow(10, 9))
	// Result: 1
	x := c / (math.Pow(
		float64(ByteUnitMap[targetUnit].base),
		float64(ByteUnitMap[targetUnit].exponential),
	))

	return x
}

func (b *Bytesize) ToKilobyte() float64  { return b.calc("kB") }
func (b *Bytesize) ToMegabyte() float64  { return b.calc("MB") }
func (b *Bytesize) ToGigabyte() float64  { return b.calc("GB") }
func (b *Bytesize) ToTerabyte() float64  { return b.calc("TB") }
func (b *Bytesize) ToPetabyte() float64  { return b.calc("PB") }
func (b *Bytesize) ToExabyte() float64   { return b.calc("EB") }
func (b *Bytesize) ToZettabyte() float64 { return b.calc("ZB") }
func (b *Bytesize) ToYottabyte() float64 { return b.calc("YB") }

func (b *Bytesize) ToKibibyte() float64 { return b.calc("KiB") }
func (b *Bytesize) ToMebibyte() float64 { return b.calc("MiB") }
func (b *Bytesize) ToGibibyte() float64 { return b.calc("GiB") }
func (b *Bytesize) ToTebibyte() float64 { return b.calc("TiB") }
func (b *Bytesize) ToPebibyte() float64 { return b.calc("PiB") }
func (b *Bytesize) ToExbibyte() float64 { return b.calc("EiB") }
func (b *Bytesize) ToZebibyte() float64 { return b.calc("ZiB") }
func (b *Bytesize) ToYobibyte() float64 { return b.calc("YiB") }

func (b *Bytesize) ToHumanReadable() string {
	// clean units
	b.Unit = b.cleanUnits(b.Unit)

	// convert given values to bytes
	c := b.Data * math.Pow(
		float64(ByteUnitMap[b.Unit].base),
		float64(ByteUnitMap[b.Unit].exponential),
	)

	// calc logarithm
	log10 := math.Log10(c)
	log10tolerant := log10 - 2

	// search ByteUnitMap for the right exponential
	var newUnit = b.Unit
	// TODO: limit to SI units?
	for key, value := range ByteUnitMap {
		if float64(value.exponential) <= log10 && float64(value.exponential) >= log10tolerant {
			newUnit = key
			break
		}
	}

	// re-calculate
	return fmt.Sprintf("%g %s", b.calc(newUnit), newUnit)
}
