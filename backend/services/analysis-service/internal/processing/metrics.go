package processing

import (
	"math"
	"sort"
)

func rmsLateral(bar []BarSample, i0, i1 int) float64 {
	if i1 <= i0 {
		return 0
	}
	x0 := bar[i0].X
	sum := 0.0
	n := 0
	for i := i0; i <= i1; i++ {
		dx := (bar[i].X - x0) * bar[i].MPP * 100.0
		sum += dx * dx
		n++
	}
	return math.Sqrt(sum / math.Max(1, float64(n)))
}

func driftX(bar []BarSample, i0, i1 int) float64 {
	if i1 <= i0 {
		return 0
	}
	dx := (bar[i1].X - bar[i0].X) * bar[i0].MPP * 100.0
	return dx
}

func jCurveDX(bar []BarSample, i0, i1 int) float64 {
	if i1 <= i0 {
		return 0
	}
	minX, maxX := bar[i0].X, bar[i0].X
	for i := i0; i <= i1; i++ {
		if bar[i].X < minX {
			minX = bar[i].X
		}
		if bar[i].X > maxX {
			maxX = bar[i].X
		}
	}
	return (maxX - minX) * bar[i0].MPP * 100.0
}

func stallCount(bar []BarSample, i0, i1 int, eps, minHoldS float64) int {
	if i1 <= i0 {
		return 0
	}

	fps := 0.0
	if i1 > i0 {
		dt := bar[i1].T - bar[i0].T
		if dt > 0 {
			fps = float64(i1-i0) / dt
		}
	}
	minLen := int(math.Ceil(minHoldS * math.Max(1, fps)))
	if minLen < 1 {
		minLen = 1
	}

	cnt := 0
	run := 0
	for i := i0; i <= i1; i++ {
		v := bar[i].VyS
		if v == 0 {
			v = bar[i].Vy
		}
		if math.Abs(v) <= eps {
			run++
		} else {
			if run >= minLen {
				cnt++
			}
			run = 0
		}
	}
	if run >= minLen {
		cnt++
	}
	return cnt
}

func eccentricP95Vy(bar []BarSample, i0, i1 int) float64 {
	if i1 <= i0 {
		return 0
	}

	var vals []float64
	for i := i0; i <= i1; i++ {
		v := bar[i].Vy
		if v == 0 {
			v = bar[i].VyS
		}
		vals = append(vals, math.Abs(v))
	}
	sort.Float64s(vals)
	k := int(0.95*float64(len(vals))) - 1
	if k < 0 {
		k = 0
	}
	return vals[k]
}

func hipsShootUpRatio(bar []BarSample, pose []PoseSample, rc RepCut, winS float64) float64 {
	if len(pose) == 0 || rc.Top <= rc.Bottom {
		return 0
	}
	startT := bar[rc.Bottom].T
	endT := math.Min(bar[rc.Top].T, startT+winS)

	pStart, okS := nearestPoseByTime(pose, startT, 30)
	pEnd, okE := nearestPoseByTime(pose, endT, 30)
	if !okS || !okE {
		return 0
	}
	_, hyS, okHS := mid2(pStart.LH, pStart.RH)
	_, hyE, okHE := mid2(pEnd.LH, pEnd.RH)
	if !(okHS && okHE) {
		return 0
	}

	barUpCm := (bar[rc.Bottom].Y - bar[rc.Bottom].Y) * bar[rc.Bottom].MPP * 100
	_ = barUpCm
	barUpCm = (bar[rc.Bottom].Y - valueAtTimeY(bar, endT)) * bar[rc.Bottom].MPP * 100.0

	hipUpCm := (hyS - hyE) * bar[rc.Bottom].MPP * 100.0
	if barUpCm <= 0 {
		return 0
	}
	return hipUpCm / barUpCm
}

func valueAtTimeY(bar []BarSample, t float64) float64 {
	if len(bar) == 0 {
		return 0
	}

	iBest := 0
	bestDT := math.Abs(bar[0].T - t)
	for i := 1; i < len(bar); i++ {
		dt := math.Abs(bar[i].T - t)
		if dt < bestDT {
			bestDT = dt
			iBest = i
		}
	}
	return bar[iBest].Y
}

func cleanBarpath(bar []BarSample) []BarSample {
	if len(bar) == 0 {
		return bar
	}
	ys := make([]float64, len(bar))
	for i := range bar {
		ys[i] = bar[i].Y
	}
	ys = movingAvg(ys, 5)
	for i := range bar {
		bar[i].Y = ys[i]
		if bar[i].MPP <= 0 || bar[i].MPP > 0.02 {
			bar[i].MPP = 0.005
		}
	}
	return bar
}

func cleanPose(pose []PoseSample) []PoseSample { return pose }
