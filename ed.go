package main

import (
	"fmt"
)

type Model struct {
	x, y string
	d    [][]int
}

func New(x, y string) *Model {
	m := new(Model)
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
	return m
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

func (m *Model) ComputeDistance() {
	for i := 1; i <= len(m.y); i++ {
		for j := 1; j <= len(m.x); j++ {
			if m.x[j-1] == m.y[i-1] {
				m.d[i][j] = m.d[i-1][j-1]
			} else {
				if m.d[i-1][j-1] < m.d[i-1][j] && m.d[i-1][j-1] < m.d[i][j-1] {
					m.d[i][j] = m.d[i-1][j-1] + 1
				} else if m.d[i-1][j] < m.d[i-1][j-1] && m.d[i-1][j] < m.d[i][j-1] {
					m.d[i][j] = m.d[i-1][j-1] + 1
				} else {
					m.d[i][j] = m.d[i][j-1] + 1
				}
			}
		}
	}
}

func main() {
	x := "CATTAG"
	y := "CAG"
	m := New(x, y)
	m.ComputeDistance()
	m.Print()
}
