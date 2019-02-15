package atkinsdiet

import (
	"reflect"
	"testing"
)

func TestSpinNotFails(t *testing.T) {
	m := New()
	_, _, err := m.Spin(1, 20)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSpin(b *testing.B) {
	m := New()
	for n := 0; n < b.N; n++ {
		m.Spin(1, 20)
	}
}

func Test_mapSymbolsToStops(t *testing.T) {
	type args struct {
		stops [3][5]int8
	}
	tests := []struct {
		name string
		args args
		want [3][5]symbol
	}{
		{
			"all zeroes",
			args{stops: [3][5]int8{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			}},
			[3][5]symbol{
				{symScale, symMayonnaise, symHam, symHam, symBacon},
				{symScale, symMayonnaise, symHam, symHam, symBacon},
				{symScale, symMayonnaise, symHam, symHam, symBacon},
			},
		},
		{
			"0-0-0-0-0",
			args{stops: [3][5]int8{
				{31, 31, 31, 31, 31},
				{0, 0, 0, 0, 0},
				{1, 1, 1, 1, 1},
			}},
			[3][5]symbol{
				{symButter, symBacon, symSteak, symBacon, symBuffaloWings},
				{symScale, symMayonnaise, symHam, symHam, symBacon},
				{symMayonnaise, symBuffaloWings, symButter, symCheese, symScale},
			},
		},
		{
			"0-1-2-3-4",
			args{stops: [3][5]int8{
				{31, 0, 1, 2, 3},
				{0, 1, 2, 3, 4},
				{1, 2, 3, 4, 5},
			}},
			[3][5]symbol{
				{symButter, symMayonnaise, symButter, symAtkins, symHam},
				{symScale, symBuffaloWings, symEggs, symScale, symCheese},
				{symMayonnaise, symSteak, symScale, symButter, symSausage},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapSymbolsToStops(tt.args.stops); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapIconsToStops() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkLine(t *testing.T) {
	type args struct {
		lineSymbols [5]symbol
		n           int8
		sym         symbol
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"3 whilds in w-w-w-1-1",
			args{
				[5]symbol{symAtkins, symAtkins, symAtkins, symHam, symHam},
				3,
				symAtkins,
			},
			true,
		},
		{
			"5 1s in w-w-w-1-1",
			args{
				[5]symbol{symAtkins, symAtkins, symAtkins, symHam, symHam},
				5,
				symHam,
			},
			true,
		},
		{
			"5 1s in w-w-w-1-2",
			args{
				[5]symbol{symAtkins, symAtkins, symAtkins, symHam, symBacon},
				4,
				symHam,
			},
			true,
		},
		{
			"3 ws in w-w-w-1-2",
			args{
				[5]symbol{symAtkins, symAtkins, symAtkins, symHam, symBacon},
				3,
				symAtkins,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkLine(tt.args.lineSymbols, tt.args.n, tt.args.sym); got != tt.want {
				t.Errorf("checkLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_caclculateLinesWin(t *testing.T) {
	type args struct {
		symbols [3][5]symbol
		bet     int
		lines   int
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"",
			args{
				[3][5]symbol{
					{symButter, symEggs, symSausage, symMayonnaise, symScale},
					{symCheese, symCheese, symBacon, symBacon, symSteak},
					{symEggs, symBacon, symCheese, symBuffaloWings, symHam},
				},
				1,
				20,
			},
			15,
		},
		{
			"",
			args{
				[3][5]symbol{
					{symBacon, symBacon, symEggs, symSausage, symScale},
					{symEggs, symSteak, symBacon, symCheese, symSteak},
					{symCheese, symSausage, symMayonnaise, symEggs, symHam},
				},
				2,
				20,
			},
			20,
		},
		{
			"",
			args{
				[3][5]symbol{
					{symAtkins, symCheese, symBuffaloWings, symCheese, symAtkins},
					{symBacon, symBacon, symBacon, symButter, symButter},
					{symMayonnaise, symEggs, symButter, symHam, symBuffaloWings},
				},
				2,
				20,
			},
			40,
		},
		{
			"",
			args{
				[3][5]symbol{
					{symMayonnaise, symHam, symBacon, symMayonnaise, symButter},
					{symHam, symCheese, symSteak, symSausage, symBuffaloWings},
					{symSausage, symBacon, symBuffaloWings, symCheese, symMayonnaise},
				},
				2,
				20,
			},
			8,
		},
		{
			"",
			args{
				[3][5]symbol{
					{symMayonnaise, symMayonnaise, symButter, symButter, symMayonnaise},
					{symHam, symButter, symMayonnaise, symHam, symEggs},
					{symCheese, symCheese, symCheese, symMayonnaise, symHam},
				},
				2,
				20,
			},
			80,
		},
		{
			"",
			args{
				[3][5]symbol{
					{symSausage, symButter, symEggs, symCheese, symSausage},
					{symBuffaloWings, symCheese, symAtkins, symAtkins, symHam},
					{symSteak, symBacon, symBuffaloWings, symScale, symAtkins},
				},
				5,
				20,
			},
			0,
		},
		{
			"",
			args{
				[3][5]symbol{
					{symSausage, symButter, symEggs, symCheese, symSausage},
					{symAtkins, symAtkins, symSteak, symHam, symBuffaloWings},
					{symSteak, symBacon, symBuffaloWings, symScale, symAtkins},
				},
				1,
				1,
			},
			40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := caclculateLinesWin(tt.args.symbols, tt.args.bet, tt.args.lines); got != tt.want {
				t.Errorf("caclculateLinesWin() = %v, want %v", got, tt.want)
			}
		})
	}
}
