package atkinsdiet

type symbol int8

const (
	// symAtkins is a wild
	symAtkins       = symbol(1)
	symSteak        = symbol(2)
	symHam          = symbol(3)
	symBuffaloWings = symbol(4)
	symSausage      = symbol(5)
	symEggs         = symbol(6)
	symButter       = symbol(7)
	symCheese       = symbol(8)
	symBacon        = symbol(9)
	symMayonnaise   = symbol(10)
	symScale        = symbol(11)

	freeSpinsAmount     = 10
	freeSpinsMultiplier = 3
)

var paytable = [34]struct {
	win    int16
	n      int8
	symbol symbol
}{
	{5000, 5, symAtkins},
	{1000, 5, symSteak},
	{500, 4, symAtkins},
	{500, 5, symHam},
	{300, 5, symBuffaloWings},
	{200, 4, symSteak},
	{200, 5, symSausage},
	{200, 5, symEggs},
	{150, 4, symHam},
	{100, 4, symBuffaloWings},
	{100, 5, symButter},
	{100, 5, symCheese},
	{75, 4, symSausage},
	{75, 4, symEggs},
	{50, 3, symAtkins},
	{50, 4, symButter},
	{50, 4, symCheese},
	{50, 5, symBacon},
	{50, 5, symMayonnaise},
	{40, 3, symSteak},
	{30, 3, symHam},
	{25, 3, symBuffaloWings},
	{25, 4, symBacon},
	{25, 4, symMayonnaise},
	{20, 3, symSausage},
	{20, 3, symEggs},
	{15, 3, symButter},
	{15, 3, symCheese},
	{10, 3, symBacon},
	{10, 3, symMayonnaise},
	{2, 5, symAtkins},
	{3, 2, symSteak},
	{2, 2, symHam},
	{2, 2, symBuffaloWings},
}

var reelStrips = [5][32]symbol{
	{symScale, symMayonnaise, symHam, symSausage, symBacon, symEggs, symCheese, symMayonnaise, symSausage, symButter, symBuffaloWings, symBacon, symEggs, symMayonnaise, symSteak, symBuffaloWings, symButter, symCheese, symEggs, symAtkins, symBacon, symMayonnaise, symHam, symCheese, symEggs, symScale, symButter, symBacon, symSausage, symBuffaloWings, symSteak, symButter},
	{symMayonnaise, symBuffaloWings, symSteak, symSausage, symCheese, symMayonnaise, symHam, symButter, symBacon, symSteak, symSausage, symMayonnaise, symHam, symAtkins, symButter, symEggs, symCheese, symBacon, symSausage, symBuffaloWings, symScale, symMayonnaise, symButter, symCheese, symBacon, symEggs, symBuffaloWings, symMayonnaise, symSteak, symHam, symCheese, symBacon},
	{symHam, symButter, symEggs, symScale, symCheese, symMayonnaise, symButter, symHam, symSausage, symBacon, symSteak, symBuffaloWings, symButter, symMayonnaise, symCheese, symSausage, symEggs, symBacon, symMayonnaise, symBuffaloWings, symHam, symSausage, symBacon, symCheese, symEggs, symAtkins, symBuffaloWings, symBacon, symButter, symCheese, symMayonnaise, symSteak},
	{symHam, symCheese, symAtkins, symScale, symButter, symBacon, symCheese, symSausage, symSteak, symEggs, symBacon, symMayonnaise, symSausage, symCheese, symButter, symHam, symMayonnaise, symBacon, symBuffaloWings, symSausage, symCheese, symEggs, symButter, symBuffaloWings, symBacon, symMayonnaise, symEggs, symHam, symSausage, symSteak, symMayonnaise, symBacon},
	{symBacon, symScale, symSteak, symHam, symCheese, symSausage, symButter, symBacon, symBuffaloWings, symCheese, symSausage, symHam, symButter, symSteak, symMayonnaise, symEggs, symSausage, symHam, symAtkins, symButter, symBuffaloWings, symMayonnaise, symEggs, symHam, symBacon, symButter, symSteak, symMayonnaise, symSausage, symEggs, symCheese, symBuffaloWings},
}

var paylines = [20][5]int8{
	{1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0},
	{2, 2, 2, 2, 2},
	{0, 1, 2, 1, 0},
	{2, 1, 0, 1, 2},
	{1, 0, 0, 0, 1},
	{1, 2, 2, 2, 1},
	{0, 0, 1, 2, 2},
	{2, 2, 1, 0, 0},
	{1, 0, 1, 2, 1},
	{1, 2, 1, 0, 1},
	{0, 1, 1, 1, 0},
	{2, 1, 1, 1, 2},
	{0, 1, 0, 1, 0},
	{2, 1, 2, 1, 2},
	{1, 1, 0, 1, 1},
	{1, 1, 2, 1, 1},
	{0, 0, 2, 0, 0},
	{2, 2, 0, 2, 2},
	{0, 2, 2, 2, 0},
}
