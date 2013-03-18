// Copyright 2013 Kevin Gillette. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package frac

import "testing"

var factors = []struct {
	a, b, gcd, lcm int
}{
	{2, 3, 1, 6},
	{2, 4, 2, 4},
	{3, 5, 1, 15},
	{-3, 5, 1, 15},
	{3, -5, 1, 15},
	{-3, -5, 1, 15},
}

func TestFactors(t *testing.T) {
	for i, f := range factors {
		if gcd := GCD(f.a, f.b); gcd != f.gcd {
			t.Errorf("GCD(%d, %d) => %d != %d [#%d]", f.a, f.b, gcd, f.gcd, i)
		}
		if lcm := LCM(f.a, f.b); lcm != f.lcm {
			t.Errorf("LCM(%d, %d) => %d != %d [#%d]", f.a, f.b, lcm, f.lcm, i)
		}
	}
}

var fracs = []struct {
	a, b, na, nb, nt, r, add, sub, mul, div F
}{
	{F{1, 2}, F{1, 2}, F{1, 2}, F{1, 2}, F{1, 2}, F{1, 2}, F{1, 1}, F{0, 1}, F{1, 4}, F{1, 1}},
	{F{2, 1}, F{2, 1}, F{2, 1}, F{2, 1}, F{2, 1}, F{2, 1}, F{4, 1}, F{0, 1}, F{4, 1}, F{1, 1}},
	{F{1, 2}, F{2, 1}, F{1, 2}, F{4, 2}, F{1, 2}, F{1, 2}, F{5, 2}, F{-3, 2}, F{1, 1}, F{1, 4}},

	{F{12, 2}, F{3, 4}, F{24, 4}, F{3, 4}, F{24, 4}, F{6, 1}, F{27, 4}, F{21, 4}, F{9, 2}, F{8, 1}},
	{F{4, 3}, F{3, 2}, F{8, 6}, F{9, 6}, F{8, 6}, F{4, 3}, F{17, 6}, F{-1, 6}, F{2, 1}, F{8, 9}},
}

func TestFracs(t *testing.T) {
	for i, f := range fracs {
		if a, b := Norm(f.a, f.b); a != f.na || b != f.nb {
			t.Errorf("Norm(%+v, %+v) => (%+v %+v) != (%+v %+v) [#%d]", f.a, f.b, a, b, f.na, f.nb, i)
		}
		if r := f.a.Reduce(); r != f.r {
			t.Fatalf("%+v.Reduce() => %+v != %+v [#%d]", f.a, r, f.r, i)
		}
		if r := f.a.NormTo(f.b); r != f.nt {
			t.Errorf("%+v.NormTo(%+v) => %+v != %+v [#%d]", f.a, f.b, r, f.nt, i)
		}
		if r := f.a.Add(f.b).Reduce(); r != f.add {
			t.Errorf("%+v + %+v => %+v != %+v [#%d]", f.a, f.b, r, f.add, i)
		}
		if r := f.a.Sub(f.b).Reduce(); r != f.sub {
			t.Errorf("%+v - %+v => %+v != %+v [#%d]", f.a, f.b, r, f.sub, i)
		}
		if r := f.a.Mul(f.b).Reduce(); r != f.mul {
			t.Errorf("%+v * %+v => %+v != %+v [#%d]", f.a, f.b, r, f.mul, i)
		}
		if r := f.a.Div(f.b).Reduce(); r != f.div {
			t.Errorf("%+v / %+v => %+v != %+v [#%d]", f.a, f.b, r, f.div, i)
		}
	}
}

var strs = []struct {
	i, o string
	f    F
}{
	{"3/2", "3/2", F{3, 2}},
	{"3/-2", "-3/2", F{3, -2}},
	{"-3/-2", "3/2", F{-3, -2}},
}

func TestStr(t *testing.T) {
	for i, c := range strs {
		if f, err := Parse(c.i); err != nil {
			t.Errorf("Parse error [#%d]: %s", i, err)
		} else if f != c.f {
			t.Errorf("%+v != %+v [#%d]", f, c.f, i)
		} else if f.String() != c.o {
			t.Errorf("%q != %q [#%d]", f.String(), c.o, i)
		}
	}
}

func BenchmarkGCD(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		GCD(i, (i+199)*3%211)
	}
}

func BenchmarkLCM(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		LCM(i, (i+199)*3%211)
	}
}

func BenchmarkNormLike(b *testing.B) {
	x := F{2, 1}
	for i := 0; i < b.N; i++ {
		x.D++
		Norm(x, F{i, x.D})
	}
}

func BenchmarkNormUnlike(b *testing.B) {
	x := F{2, 1}
	for i := 0; i < b.N; i++ {
		x.D++
		Norm(x, F{x.D, x.N})
	}
}

func BenchmarkNormToLike(b *testing.B) {
	x := F{2, 1}
	for i := 0; i < b.N; i++ {
		x.D++
		x.NormTo(F{i, x.D})
	}
}

func BenchmarkNormToUnlike(b *testing.B) {
	x := F{2, 1}
	for i := 0; i < b.N; i++ {
		x.D++
		x.NormTo(F{x.D, x.N})
	}
}

func BenchmarkReduce(b *testing.B) {
	a := F{2, 2}
	for i := 0; i < b.N; i++ {
		a.N += i % 10
		a.Reduce()
		a = F{a.D, a.N}
	}
}

func BenchmarkSubLike(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		F{3, 17}.Sub(F{i, 17})
	}
}

func BenchmarkSubUnlike(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		F{3, 17}.Sub(F{1, i})
	}
}

func BenchmarkDiv(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		F{3, 17}.Div(F{i, 23})
	}
}

func BenchmarkInv(b *testing.B) {
	offset := b.N / 2
	for i := -offset; i <= offset; i++ {
		F{i, 17}.Inv()
	}
}
