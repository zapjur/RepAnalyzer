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
	if n < 3 {
		return nil
	}

	y := make([]float64, n)
	for i := range bar {
		y[i] = bar[i].Y
	}
	yS := movingAvg(y, 7)

	vyStart := thresholds["seg_vy_start"]
	vyStop := thresholds["seg_vy_stop"]
	minDur := thresholds["seg_min_dur_s"]
	minRise := thresholds["seg_min_vert_cm"]
	quietS := thresholds["seg_quiet_s"]

	type stateT int
	const (
		idle stateT = iota
		up
	)

	state := idle
	lastTop := -1
	quietSince := -1.0
	bottomIdx := -1

	// pomocnicze: znajdź max/min yS w oknie czasowym wstecz (np. 0.25 s)
	findLocalMaxBack := func(i int, winS float64) int {
		best := i
		tNow := bar[i].T
		for j := i; j >= 0 && tNow-bar[j].T <= winS; j-- {
			if yS[j] > yS[best] {
				best = j
			}
		}
		return best
	}
	findLocalMinBack := func(i int, winS float64) int {
		best := i
		tNow := bar[i].T
		for j := i; j >= 0 && tNow-bar[j].T <= winS; j-- {
			if yS[j] < yS[best] {
				best = j
			}
		}
		return best
	}

	var reps []RepCut

	for i := 0; i < n; i++ {
		v := bar[i].VyS
		if v == 0 {
			v = bar[i].Vy
		}

		switch state {
		case idle:
			if v > vyStart {
				bottomIdx = findLocalMaxBack(i, 0.25)
				state = up
				quietSince = -1
			}
		case up:
			if v < vyStop {
				if quietSince < 0 {
					quietSince = bar[i].T
				}
				if bar[i].T-quietSince >= quietS {
					topIdx := findLocalMinBack(i, 0.25)
					if bottomIdx >= 0 && topIdx > bottomIdx {
						dt := bar[topIdx].T - bar[bottomIdx].T
						mpp := bar[bottomIdx].MPP
						if mpp <= 0 {
							mpp = bar[topIdx].MPP
						}
						dyCm := (yS[bottomIdx] - yS[topIdx]) * mpp * 100.0
						if dt >= minDur && dyCm >= minRise {
							reps = append(reps, RepCut{Bottom: bottomIdx, Top: topIdx, PrevTop: lastTop})
							lastTop = topIdx
						}
					}
					state = idle
					quietSince = -1
					bottomIdx = -1
				}
			} else {
				quietSince = -1
			}
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
