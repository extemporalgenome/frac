// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package frac provides a lightweight rational type.
//
// Emphasis is on simplicity, reliability, and speed, and unlike many
// rational type systems, frac exposes internals, does not auto-simplify, and
// has no limiting invariants (the fields can take any combination of values).
// The design is oriented toward practical use of ratios and fractions,
// particularly multimedia (such as video sizes), where it is often important
// to work with data as either a simplified ratio or a non-simplified pair of
// integers, depending on the context.
//
// The sole type is F, and the numerator and denominator are exported through
// the fields N and D. frac does not provide trivial operations since those
// can be done explicitly using the exported fields. All methods operate on
// value receivers, and mutating methods are chainable.
//
// For all arithemtic operations over x and y, of type F, the resulting
// denominator will never be smaller than max(x.D, y.D). Use Reduce or NormTo
// to refine the result or prevent overflow as needed.
package frac

import "strconv"

func Parse(s string) (F, error) {
	var err error
	t := ""
	x := F{0, 1}
	for i, r := range s {
		if r == '/' {
			s, t = s[:i], s[i+1:]
			break
		}
	}
	x.N, err = strconv.Atoi(s)
	if err == nil && t != "" {
		x.D, err = strconv.Atoi(t)
	}
	return x, err
}

// F represents a fraction/ratio/rational. There are no invariants, so all its
// methods will complete without an error or panic even if either or both fields
// are negative or if D is zero. To simulate a "divide by zero" panic, see
// Assert.  All arithemtic methods sacrifice overflow protection for speed.
type F struct {
	N, D int
}

// Mul returns the product of x*y.
func (x F) Mul(y F) F {
	return F{x.N * y.N, x.D * y.D}
}

// Div returns the quotient of x/y.
func (x F) Div(y F) F {
	if y.N < 0 {
		y.N, y.D = -y.N, -y.D
	}
	return F{x.N * y.D, x.D * y.N}
}

// Add returns the sum of x+y.
func (x F) Add(y F) F {
	if x.D != y.D {
		x.N *= y.D
		y.N *= x.D
		x.D *= y.D
	}
	return F{x.N + y.N, x.D}
}

// Sub returns the difference of x-y.
func (x F) Sub(y F) F {
	if x.D != y.D {
		x.N *= y.D
		y.N *= x.D
		x.D *= y.D
	}
	return F{x.N - y.N, x.D}
}

// Reduce returns an equivalent to x represented in lowest terms.
func (x F) Reduce() F {
	d := GCD(x.N, x.D)
	x.N /= d
	x.D /= d
	if x.D < 0 {
		x.N, x.D = -x.N, -x.D
	}
	return x
}

// NormTo returns an equivalent to x scaled to match the minimum shared
// multiple of y.D.
func (x F) NormTo(y F) F {
	if x.D == y.D {
		return x
	}
	m := LCM(x.Reduce().D, y.D)
	return F{m / x.D * x.N, m}
}

// Inv returns the inversion of x.
func (x F) Inv() F {
	if x.N < 0 {
		return F{-x.D, -x.N}
	}
	return F{x.D, x.N}
}

// Abs returns the rational distance of x from zero.
func (x F) Abs() F {
	if x.N < 0 {
		x.N = -x.N
	}
	if x.D < 0 {
		x.D = -x.D
	}
	return x
}

// Cmp comparse x and y and returns:
//     -1 if x <  y
//      0 if x == y
//     +1 if x >  y
func (x F) Cmp(y F) int {
	if x.D != y.D {
		x.N *= y.D
		y.N *= x.D
	}
	switch {
	case x.N < y.N:
		return -1
	case x.N > y.N:
		return +1
	}
	return 0
}

// IsNeg returns true if F is arithmetically negative, otherwise false.
func (x F) IsNeg() bool { return x.N^x.D < 0 }

// Float64 returns the floating-point approximation of x.
func (x F) Float64() float64 { return float64(x.N) / float64(x.D) }

// Assert panics if x.D is zero, returning x otherwise.
func (x F) Assert() F {
	_ = 1 / x.D
	return x
}

func (x F) String() string {
	b := make([]byte, 0, 40)
	if x.D < 0 {
		x.N = -x.N
		x.D = -x.D
	}
	b = strconv.AppendInt(b, int64(x.N), 10)
	b = append(b, '/')
	b = strconv.AppendInt(b, int64(x.D), 10)
	return string(b)
}

// Norm returns bidirectionally normalized equivalents to x and y. Unlike
// the NormTo method, the resulting denominators may be smaller than those of
// either input.
func Norm(x, y F) (F, F) {
	if x.D == y.D {
		return x, y
	}
	m := LCM(x.Reduce().D, y.Reduce().D)
	if m <= x.D {
		m = x.D
	} else if m < y.D {
		m = y.D
	}
	return F{m / x.D * x.N, m}, F{m / y.D * y.N, m}
}

// GCD returns the greatest common divisor of x and y. The result may be
// negative if either argument is negative.
func GCD(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return abs(x)
}

// LCM returns the least common multiple of x and y. The result may be
// negative if either argument is negative.
func LCM(x, y int) int {
	if x == y || x == -y {
		return x
	}
	return abs(x / GCD(x, y) * y)
}

// abs returns the distance of the argument from zero.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
