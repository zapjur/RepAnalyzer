package processing

func movingAvg(xs []float64, k int) []float64 {
	if k < 3 {
		k = 3
	}
	if k%2 == 0 {
		k++
	}
	n := len(xs)
	if n == 0 {
		return nil
	}
	out := make([]float64, n)
	half := k / 2
	sum := 0.0
	cnt := 0
	for i := 0; i < n; i++ {
		lo := maxInt(0, i-half)
		hi := minInt(n-1, i+half)
		if i == 0 {
			sum, cnt = 0, 0
			for j := lo; j <= hi; j++ {
				sum += xs[j]
				cnt++
			}
		} else {
			prevLo := maxInt(0, i-1-half)
			prevHi := minInt(n-1, i-1+half)
			if lo > prevLo {
				sum -= xs[prevLo]
				cnt--
			}
			if hi > prevHi {
				sum += xs[hi]
				cnt++
			}
		}
		if cnt > 0 {
			out[i] = sum / float64(cnt)
		} else {
			out[i] = xs[i]
		}
	}
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func segmentConcentric(bar []BarSample) []RepCut {
	n := len(bar)
	if n < 5 {
		return nil
	}
	y := make([]float64, n)
	for i := range bar {
		y[i] = bar[i].Y
	}
	yS := movingAvg(y, 7)

	isMax := func(i int) bool {
		if i <= 0 || i >= n-1 {
			return false
		}
		return yS[i] >= yS[i-1] && yS[i] > yS[i+1]
	}
	isMin := func(i int) bool {
		if i <= 0 || i >= n-1 {
			return false
		}
		return yS[i] <= yS[i-1] && yS[i] < yS[i+1]
	}

	var bottoms, tops []int
	for i := 1; i < n-1; i++ {
		if isMax(i) {
			bottoms = append(bottoms, i)
		}
		if isMin(i) {
			tops = append(tops, i)
		}
	}
	if len(bottoms) == 0 || len(tops) == 0 {
		return nil
	}

	var reps []RepCut
	ti := 0
	lastTop := 0
	for _, b := range bottoms {
		for ti < len(tops) && tops[ti] <= b {
			lastTop = tops[ti]
			ti++
		}
		if ti >= len(tops) {
			break
		}
		t := tops[ti]
		if t <= b {
			continue
		}
		dt := bar[t].T - bar[b].T
		dyPx := yS[b] - yS[t]
		mpp := bar[b].MPP
		if mpp <= 0 {
			mpp = bar[t].MPP
		}
		dyCm := dyPx * mpp * 100.0
		if dt >= thresholds["seg_min_dur_s"] && dyCm >= thresholds["seg_min_vert_cm"] {
			reps = append(reps, RepCut{Bottom: b, Top: t, PrevTop: lastTop})
		}
	}
	return reps
}

func mergeReps(bar []BarSample, reps []RepCut, fps float64, minDescentCm, maxGapS float64) []RepCut {
	if len(reps) <= 1 {
		return reps
	}
	var out []RepCut
	cur := reps[0]
	for i := 1; i < len(reps); i++ {
		nxt := reps[i]
		gap := bar[nxt.Bottom].T - bar[cur.Top].T
		yTop := bar[cur.Top].Y
		yNextBottom := bar[nxt.Bottom].Y
		descentCm := (yNextBottom - yTop) * bar[cur.Top].MPP * 100.0
		if gap <= maxGapS && descentCm < minDescentCm {
			cur.Top = nxt.Top
		} else {
			out = append(out, cur)
			cur = nxt
		}
	}
	out = append(out, cur)
	return out
}
