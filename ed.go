// Author: Vinhthuy Phan,  June 2015
// Find all alignments
package main

import (
	"errors"
	"fmt"
)

type Variant struct {
	ref string
	alt string
}

func NewVariant(ref, alt string) *Variant {
	return &Variant{ref, alt}
}

func (v *Variant) Print() {
	fmt.Println(v.ref, "|", v.alt)
}

// Alignment
type Alignment struct {
	X, Y string
}

func (a *Alignment) Print() {
	fmt.Println(a.X)
	fmt.Println(a.Y)
}

// Call variants for this alignment.  Do not call potentially incomplete variants (PIV).
// PIV occur at the beginning or end of the alignment.
func (a *Alignment) Call() ([]*Variant, bool, bool, error) {
	PIVb, PIVe := false, false
	n := len(a.X) - 1 // X, Y have the same length
	b, e := 0, n
	variants := make([]*Variant, 0)

	// Skip PIV
	if a.X[0] == '-' || a.Y[0] == '-' {
		PIVb = true
		for ; b <= n; b++ {
			if a.X[b] == '-' && a.Y[b] == '-' {
				return variants, PIVb, PIVe, errors.New("Alginment error: two gaps aligned.")
			}
			if a.X[b] != '-' && a.Y[b] != '-' {
				break
			}
		}
	}
	if a.X[n] == '-' || a.Y[n] == '-' {
		PIVe = true
		for ; e >= 0; e-- {
			if a.X[e] == '-' && a.Y[e] == '-' {
				return variants, PIVb, PIVe, errors.New("Alginment error: two gaps aligned.")
			}
			if a.X[e] != '-' && a.Y[e] != '-' {
				break
			}
		}
	}

	cur_var := ""

	for i := b; i <= e; i++ {
		if a.X[i] == '-' && a.Y[i] == '-' {
			return variants, PIVb, PIVe, errors.New("Alignment error: two gaps aligned.")
		}
		if a.X[i] != a.Y[i] {
			if a.X[i] != '-' && a.Y[i] != '-' {
				// SNP
				variants = append(variants, NewVariant(string(a.X[i]), string(a.Y[i])))
			} else if a.X[i] == '-' {
				if a.X[i-1] != '-' && a.X[i-1] != a.Y[i-1] {
					return variants, PIVb, PIVe, errors.New("Invalid indel configuration: " + "\n\t" + a.X[i-1:] + "\n\t" + a.Y[i-1:])
				}
				// insertion
				cur_var += string(a.Y[i-1])
				if a.X[i+1] != '-' {
					cur_var += string(a.Y[i])
					variants = append(variants, NewVariant(string(cur_var[0]), cur_var))
					cur_var = ""
				}
			} else if a.Y[i] == '-' {
				if a.Y[i-1] != '-' && a.Y[i-1] != a.X[i-1] {
					return variants, PIVb, PIVe, errors.New("Invalid indel configuration: " + "\n\t" + a.Y[i-1:] + "\n\t" + a.X[i-1:])
				}
				// deletion
				cur_var += string(a.X[i-1])
				if a.Y[i+1] != '-' {
					cur_var += string(a.X[i])
					variants = append(variants, NewVariant(cur_var, string(cur_var[0])))
					cur_var = ""
				}
			}
		}
	}

	return variants, PIVb, PIVe, nil
}

// Solution
type Solution struct {
	A []*Alignment
}

func NewSolution() *Solution {
	s := new(Solution)
	s.A = make([]*Alignment, 0)
	return s
}

func (s *Solution) Extend(cx, cy byte) {
	if len(s.A) == 0 {
		s.A = append(s.A, &Alignment{string(cx), string(cy)})
	} else {
		for _, a := range s.A {
			a.X += string(cx)
			a.Y += string(cy)
		}
	}
}

func (s *Solution) Merge(t *Solution) {
	if t != nil {
		s.A = append(s.A, t.A...)
	}
}

func (s *Solution) Print() {
	fmt.Println("Number of alignments:", len(s.A))
	for _, a := range s.A {
		a.Print()
	}
}

// A model of alignment
// x is reference.  y is aligned to x
type Model struct {
	x, y string
	d    [][]int
}

func New() *Model {
	return new(Model)
}

func (m *Model) Init(x, y string) {
	m.x = x
	m.y = y
	m.d = make([][]int, len(y)+1)
	for i := 0; i <= len(y); i++ {
		m.d[i] = make([]int, len(x)+1)
	}
	for i := 0; i <= len(y); i++ {
		m.d[i][0] = i
	}
	for j := 0; j <= len(x); j++ {
		m.d[0][j] = j
	}
}

func (m *Model) Print() {
	fmt.Print("     ")
	for j := 0; j < len(m.x); j++ {
		fmt.Printf("%2s ", string(m.x[j]))
	}
	fmt.Println()

	for i := 0; i <= len(m.y); i++ {
		if i > 0 {
			fmt.Print(string(m.y[i-1]), " ")
		} else {
			fmt.Print("  ")
		}
		for j := 0; j <= len(m.x); j++ {
			fmt.Printf("%2d ", m.d[i][j])
		}
		fmt.Println()
	}
}

func (m *Model) ComputeDistance(x, y string) {
	m.Init(x, y)
	var mis int
	for i := 1; i <= len(m.y); i++ {
		for j := 1; j <= len(m.x); j++ {
			if m.x[j-1] == m.y[i-1] {
				mis = 0
			} else {
				mis = 1
			}
			if m.d[i-1][j-1]+mis <= m.d[i-1][j]+1 && m.d[i-1][j-1]+mis <= m.d[i][j-1]+1 {
				m.d[i][j] = m.d[i-1][j-1] + mis
			} else if m.d[i-1][j]+1 <= m.d[i-1][j-1]+mis && m.d[i-1][j] <= m.d[i][j-1] {
				m.d[i][j] = m.d[i-1][j] + 1
			} else {
				m.d[i][j] = m.d[i][j-1] + 1
			}
		}
	}
}

func (m *Model) Trace() *Solution {
	return m.trace(len(m.y), len(m.x))
}

func (m *Model) trace(i, j int) *Solution {
	var mis int
	var Ssub, Sins, Sdel *Solution
	ret := NewSolution()

	if i > 0 && j > 0 {
		if m.x[j-1] == m.y[i-1] {
			mis = 0
		} else {
			mis = 1
		}

		if m.d[i-1][j-1]+mis <= m.d[i-1][j]+1 && m.d[i-1][j-1]+mis <= m.d[i][j-1]+1 {
			Ssub = m.trace(i-1, j-1)
			Ssub.Extend(m.x[j-1], m.y[i-1])
			ret = Ssub
		}
		if m.d[i-1][j]+1 <= m.d[i-1][j-1]+mis && m.d[i-1][j] <= m.d[i][j-1] {
			Sins = m.trace(i-1, j)
			Sins.Extend('-', m.y[i-1])
			ret.Merge(Sins)
		}
		if m.d[i][j-1]+1 <= m.d[i-1][j-1]+mis && m.d[i][j-1] <= m.d[i-1][j] {
			Sdel = m.trace(i, j-1)
			Sdel.Extend(m.x[j-1], '-')
			ret.Merge(Sdel)
		}
	} else {
		for ; i > 0; i-- {
			ret.Extend('-', m.y[i-1])
		}
		for ; j > 0; j-- {
			ret.Extend(m.x[j-1], '-')
		}
	}
	return ret
}

func (m *Model) Call(sol *Solution) {
	if sol != nil {
		for _, algn := range sol.A {
			algn.Print()
			variants, pivb, pive, err := algn.Call()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(pivb, pive)
				for _, v := range variants {
					v.Print()
				}
			}
		}
	}
}

func main() {
	// y := "CT"
	// x := "C"
	// y := "CATTAG"
	// x := "CGGTAG"
	// y := "CATTAG"
	// x := "CTTAG"
	// x := "CATTAG"
	// y := "CAG"
	// x := "CATGATG"
	// y := "CATG"
	// x := "CATCCATG"
	// y := "CATG"
	y := "CATCCATG"
	x := "GATG"
	m := New()
	m.ComputeDistance(x, y)
	m.Print()
	solution := m.Trace()
	// solution.Print()
	m.Call(solution)

	// a := Alignment{"-AC--TGGTAGGT", "AACGCTGAT-G--"}
	// variants, pivb, pive, err := a.Call()
	// fmt.Println(pivb, pive, err)
	// a.Print()
	// for _, v := range variants {
	// 	v.Print()
	// }
}
