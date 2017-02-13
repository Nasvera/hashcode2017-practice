package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Pizza struct {
	Rows          int
	Cols          int
	Min           int
	Max           int
	Data          [][]Incredient
	MushroomCount int
	TomatoCount   int
}

type Point struct {
	X int
	Y int
}

type Rect struct {
	A        Point
	B        Point
	vertical bool
}

type Node struct {
	Left   *Node
	Right  *Node
	Rect   Rect
	MC     int
	TC     int
	Leaf   bool
	Valid  bool
	Vcheck bool
	Hcheck bool
}

type Incredient byte

const (
	Tomato   Incredient = 'T'
	Mushroom Incredient = 'M'
)

const USAGE = `
    go run *.go [data_file]
`

var (
	score  int
	err    error
	f      *os.File
	pizza  = &Pizza{}
	splits = make([]Rect, 0)
)

func init() {
	if len(os.Args) != 2 {
		log.Fatal(USAGE)
	}
	if f, err = os.Open(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func main() {
	parse()
	root := &Node{
		Rect: Rect{
			A: Point{X: 0, Y: 0},
			B: Point{X: pizza.Cols - 1, Y: pizza.Rows - 1},
		},
		MC: pizza.MushroomCount,
		TC: pizza.TomatoCount,
	}
	splitVertically(root)
	output("v.out")

	score = 0
	splits = make([]Rect, 0)

	splitHorizontally(root)
	output("h.out")
}

func splitVertically(n *Node) {
	n.Vcheck = true
	// Invalid slice not enought incredients
	if n.MC < pizza.Min || n.TC < pizza.Min {
		//fmt.Println("Invalid slice not enought incredients")
		n.Leaf = true
		n.Valid = false
		return
	}

	// Check if slice is valid
	if n.MC+n.TC <= pizza.Max {
		//fmt.Println("Valid slice")
		n.Leaf = true
		n.Valid = true
		splits = append(splits, n.Rect)
		score += n.MC + n.TC
		return
	}

	// Slice can not be split vertically anymore
	if n.Rect.A.X == n.Rect.B.X {
		//fmt.Println("Slice can not be split vertically")
		if !n.Hcheck {
			splitHorizontally(n)
		} else {
			n.Leaf = true
			n.Valid = false
		}
		return
	}

	var ltc, lmc, i int
	for i = n.Rect.A.X; i <= n.Rect.B.X; i++ {
		for j := n.Rect.A.Y; j <= n.Rect.B.Y; j++ {
			if pizza.Data[j][i] == Tomato {
				ltc++
				continue
			}
			lmc++
		}

		if n.MC < n.TC {
			if lmc >= n.MC/2 {
				break
			}
		} else {
			if ltc >= n.TC/2 {
				break
			}
		}
	}
	la := n.Rect.A
	lb := Point{X: i, Y: n.Rect.B.Y}
	ra := Point{X: i + 1, Y: n.Rect.A.Y}
	rb := n.Rect.B

	if i == n.Rect.B.X {
		// Could not split slice from left to right vertically
		if !n.Hcheck {
			// Try splitting horizontaly
			splitHorizontally(n)
			if n.Valid || !n.Leaf {
				// Horizontal check was succesful
				return
			}
		}

		// Lets try reverse loop for splitting
		ltc, lmc, i = 0, 0, 0
		for i = n.Rect.B.X; i >= n.Rect.A.X; i-- {
			for j := n.Rect.A.Y; j <= n.Rect.B.Y; j++ {
				if pizza.Data[j][i] == Tomato {
					ltc++
					continue
				}
				lmc++
			}

			if n.MC < n.TC {
				if lmc >= n.MC/2 {
					break
				}
			} else {
				if ltc >= n.TC/2 {
					break
				}
			}
		}

		if i == n.Rect.A.X {
			// Reverse loop has also failed
			// Lets mark this slice as a lost case
			//fmt.Println("Invalid slice no way to split")
			n.Leaf = true
			n.Valid = false
			return
		}

		la = n.Rect.A
		lb = Point{X: i - 1, Y: n.Rect.B.Y}
		ra = Point{X: i, Y: n.Rect.A.Y}
		rb = n.Rect.B

	}

	n.Left = &Node{
		Rect: Rect{A: la, B: lb},
		MC:   lmc,
		TC:   ltc,
	}
	n.Right = &Node{
		Rect: Rect{A: ra, B: rb},
		MC:   n.MC - lmc,
		TC:   n.TC - ltc,
	}

	//splitHorizontally(n.Left)
	//splitHorizontally(n.Right)

	splitVertically(n.Left)
	splitVertically(n.Right)
}

func splitHorizontally(n *Node) {
	n.Hcheck = true
	// Invalid slice not enought incredients
	if n.MC < pizza.Min || n.TC < pizza.Min {
		//fmt.Println("Invalid slice not enought incredients")
		n.Leaf = true
		n.Valid = false
		return
	}

	// Check if slice is valid
	if n.MC+n.TC <= pizza.Max {
		//fmt.Println("Valid slice")
		n.Leaf = true
		n.Valid = true
		splits = append(splits, n.Rect)
		score += n.MC + n.TC
		return
	}

	// Slice can not be split horizontally anymore
	if n.Rect.A.X == n.Rect.B.X {
		//fmt.Println("Slice can not be split horizontally")
		if !n.Vcheck {
			splitVertically(n)
		} else {
			n.Leaf = true
			n.Valid = false
		}
		return
	}

	var ltc, lmc, i int
	for i = n.Rect.A.Y; i <= n.Rect.B.Y; i++ {
		for j := n.Rect.A.X; j <= n.Rect.B.X; j++ {
			if pizza.Data[i][j] == Tomato {
				ltc++
				continue
			}
			lmc++
		}

		if n.MC < n.TC {
			if lmc >= n.MC/2 {
				break
			}
		} else {
			if ltc >= n.TC/2 {
				break
			}
		}
	}
	la := n.Rect.A
	lb := Point{X: n.Rect.B.X, Y: i}
	ra := Point{X: n.Rect.A.X, Y: i + 1}
	rb := n.Rect.B

	if i == n.Rect.B.Y {
		// Could not split slice from left to right horizontally
		if !n.Vcheck {
			splitVertically(n)
			if n.Valid || !n.Leaf {
				// Vertical check was succesful
				return
			}

			// Lets try reverse loop for splitting
			ltc, lmc, i = 0, 0, 0
			for i = n.Rect.B.Y; i >= n.Rect.A.Y; i-- {
				for j := n.Rect.A.X; j <= n.Rect.B.X; j++ {
					if pizza.Data[i][j] == Tomato {
						ltc++
						continue
					}
					lmc++
				}

				if n.MC < n.TC {
					if lmc >= n.MC/2 {
						break
					}
				} else {
					if ltc >= n.TC/2 {
						break
					}
				}
			}

			if i == n.Rect.A.Y {
				// Reverse loop has also failed
				// Lets mark this slice as a lost case
				//fmt.Println("Invalid slice no way to split")
				n.Leaf = true
				n.Valid = false
				return
			}

			la = n.Rect.A
			lb = Point{X: n.Rect.B.X, Y: i - 1}
			ra = Point{X: n.Rect.A.X, Y: i}
			rb = n.Rect.B
		}
	}

	n.Left = &Node{
		Rect: Rect{A: la, B: lb},
		MC:   lmc,
		TC:   ltc,
	}
	n.Right = &Node{
		Rect: Rect{A: ra, B: rb},
		MC:   n.MC - lmc,
		TC:   n.TC - ltc,
	}

	//splitVertically(n.Left)
	//splitVertically(n.Right)

	splitHorizontally(n.Left)
	splitHorizontally(n.Right)
}

func parse() {
	scanner := bufio.NewScanner(f)

	// Read pizza info line
	scanner.Scan()
	params := strings.Split(scanner.Text(), " ")
	if len(params) != 4 {
		log.Fatal("Invalid lenght of parameters")
	}

	// Set parameters
	func(args ...*int) {
		for i, a := range args {
			if *a, err = strconv.Atoi(params[i]); err != nil {
				log.Fatal(err)
			}
		}
	}(&pizza.Rows, &pizza.Cols, &pizza.Min, &pizza.Max)

	// Set data for pizza
	pizza.Data = make([][]Incredient, pizza.Rows)
	rowNum := 0
	var col []byte
	for scanner.Scan() {
		if rowNum == pizza.Rows {
			log.Fatal("Too many rows")
		}

		pizza.Data[rowNum] = make([]Incredient, pizza.Cols)
		col = scanner.Bytes()
		if len(col) != pizza.Cols {
			log.Fatal("Invalid row length")
		}

		for i, b := range col {
			pizza.Data[rowNum][i] = Incredient(b)
			if pizza.Data[rowNum][i] == Tomato {
				pizza.TomatoCount++
			} else {
				pizza.MushroomCount++
			}
		}
		rowNum++
	}

	if rowNum != pizza.Rows {
		log.Fatalf("Not enaugh rows")
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func output(name string) {
	var w io.Writer
	if w, err = os.Create(name); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, len(splits))
	for _, s := range splits {
		fmt.Fprintf(w, "%d %d %d %d\n", s.A.Y, s.A.X, s.B.Y, s.B.X)
	}

	fmt.Printf("Score: %d/%d\n", score, pizza.Rows*pizza.Cols)
}
