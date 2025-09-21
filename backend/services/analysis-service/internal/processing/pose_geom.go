package processing

import (
	"math"
	"sort"
)

func sign(v float64) float64 {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	default:
		return 0
	}
}
func median(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	ys := append([]float64(nil), xs...)
	sort.Float64s(ys)
	m := len(ys) / 2
	if len(ys)%2 == 1 {
		return ys[m]
	}
	return 0.5 * (ys[m-1] + ys[m])
}
func mad(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	m := median(xs)
	dev := make([]float64, 0, len(xs))
	for _, x := range xs {
		if !math.IsNaN(x) && !math.IsInf(x, 0) {
			dev = append(dev, math.Abs(x-m))
		}
	}
	return median(dev)
}
func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mid2(a, b Point) (x, y float64, ok bool) {
	if a.Valid && b.Valid {
		return 0.5 * (a.X + b.X), 0.5 * (a.Y + b.Y), true
	}
	return 0, 0, false
}

type ankleAnchor struct {
	Side         string
	AnkleXMedian float64
	KneeXMedian  float64
	ForeSign     float64
	Samples      int
}

func chooseStableAnkle(bar []BarSample, pose []PoseSample) ankleAnchor {
	type item struct{ ax, kx, bx float64 }

	var L, R []item
	barX := map[int]float64{}
	for _, b := range bar {
		barX[b.Frame] = b.X
	}
	for _, p := range pose {
		if bx, ok := barX[p.Frame]; ok {
			if p.LA.Valid {
				kx := p.LK.X
				if !p.LK.Valid {
					kx = p.LH.X
				}
				L = append(L, item{ax: p.LA.X, kx: kx, bx: bx})
			}
			if p.RA.Valid {
				kx := p.RK.X
				if !p.RK.Valid {
					kx = p.RH.X
				}
				R = append(R, item{ax: p.RA.X, kx: kx, bx: bx})
			}
		}
	}

	score := func(S []item) (medA, medK, fore, medDist, stab float64, n int) {
		if len(S) == 0 {
			return 0, 0, 1, 1e9, 0, 0
		}
		ax, kx, dist := make([]float64, 0, len(S)), make([]float64, 0, len(S)), make([]float64, 0, len(S))
		for _, s := range S {
			ax = append(ax, s.ax)
			kx = append(kx, s.kx)
			dist = append(dist, math.Abs(s.bx-s.ax))
		}
		medA = median(ax)
		medK = median(kx)
		medDist = median(dist)
		stab = 1.0 / (mad(ax) + 1e-6)
		fore = sign(median(subVec(kx, ax)))
		if fore == 0 {
			fore = 1
		}
		return medA, medK, fore, medDist, stab, len(S)
	}

	lA, lK, lFore, lDist, lStab, lN := score(L)
	rA, rK, rFore, rDist, rStab, rN := score(R)

	useLeft := false
	switch {
	case lN > 0 && rN == 0:
		useLeft = true
	case rN > 0 && lN == 0:
		useLeft = false
	default:
		if lDist < rDist-1e-6 {
			useLeft = true
		} else if rDist < lDist-1e-6 {
			useLeft = false
		} else if lStab > rStab+1e-6 {
			useLeft = true
		} else if rStab > lStab+1e-6 {
			useLeft = false
		} else {
			useLeft = lN >= rN
		}
	}

	if useLeft {
		return ankleAnchor{"L", lA, lK, lFore, lN}
	}
	return ankleAnchor{"R", rA, rK, rFore, rN}
}

func subVec(a, b []float64) []float64 {
	n := min(len(a), len(b))
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = a[i] - b[i]
	}
	return out
}

func torsoAngleDegAt(ps PoseSample) float64 {
	hx, hy, okH := mid2(ps.LH, ps.RH)
	sx, sy, okS := mid2(ps.LS, ps.RS)
	if okH && okS {
		dx := sx - hx
		dy := hy - sy
		ang := math.Atan2(math.Abs(dy), math.Abs(dx)+1e-6) * 180 / math.Pi
		return clamp(ang, 0, 90)
	}

	if ps.LH.Valid && ps.LS.Valid {
		dx := ps.LS.X - ps.LH.X
		dy := ps.LH.Y - ps.LS.Y
		return clamp(math.Atan2(math.Abs(dy), math.Abs(dx)+1e-6)*180/math.Pi, 0, 90)
	}
	if ps.RH.Valid && ps.RS.Valid {
		dx := ps.RS.X - ps.RH.X
		dy := ps.RH.Y - ps.RS.Y
		return clamp(math.Atan2(math.Abs(dy), math.Abs(dx)+1e-6)*180/math.Pi, 0, 90)
	}
	return 0
}
