package processing

import "math"

func judgeDeadlift(r RepReport) RepReport {
	flags := []string{}
	verdict := "ok"

	if fl, ok := r.Features["bar_over_midfoot_start_fl"]; ok {
		if math.Abs(fl) > thresholds["bar_midfoot_fl_err"] {
			flags = append(flags, "bar_not_over_midfoot")
			verdict = "warn"
		} else if math.Abs(fl) > thresholds["bar_midfoot_fl_warn"] {
			flags = append(flags, "bar_slightly_off_midfoot")
			if verdict == "ok" {
				verdict = "warn"
			}
		}
	}

	if sofl, ok := r.Features["shoulder_bar_offset_start_fl"]; ok {
		if sofl > thresholds["shoulder_over_bar_fl"] {
			flags = append(flags, "shoulders_too_far_over_bar")
			verdict = "warn"
		}
	}

	if sc, ok := r.Features["stall_count"]; ok && sc >= 2 {
		flags = append(flags, "hitching")
		verdict = "warn"
	}

	if hr, ok := r.Features["hips_shoot_up_ratio"]; ok && hr > 1.5 {
		flags = append(flags, "hips_shoot_up")
		verdict = "warn"
	}

	r.Flags = flags
	r.Verdict = verdict
	return r
}

func judgeSquat(r RepReport) RepReport {
	flags := []string{}
	verdict := "ok"

	if d, ok := r.Features["depth_ok"]; ok && d < 0.5 {
		flags = append(flags, "depth_insufficient")
		verdict = "warn"
	}

	if ta, ok := r.Features["torso_angle_bottom_deg"]; ok && ta > thresholds["torso_sq_warn"] {
		flags = append(flags, "torso_lean_high")
		verdict = "warn"
	}

	if dx, ok := r.Features["drift_x_cm"]; ok && math.Abs(dx) > thresholds["drift_err"] {
		flags = append(flags, "barpath_drift")
		verdict = "warn"
	}

	r.Flags = flags
	r.Verdict = verdict
	return r
}

func judgeBench(r RepReport) RepReport {
	flags := []string{}
	verdict := "ok"

	if ev, ok := r.Features["ecc_p95_vy_m_s"]; ok {
		if thr, ok2 := thresholds["ecc_vy_warn_bench"]; ok2 && ev > thr {
			flags = append(flags, "eccentric_too_fast")
			verdict = "warn"
		}
	}

	if sc, ok := r.Features["stall_count"]; ok && sc >= 1 {
		flags = append(flags, "stall")
		verdict = "warn"
	}

	if j, ok := r.Features["jcurve_dx_cm"]; ok {
		if jmin, okMin := thresholds["jcurve_min"]; okMin && j < jmin {
			flags = append(flags, "barpath_too_linear")
			verdict = "warn"
		} else if jmax, okMax := thresholds["jcurve_max"]; okMax && j > jmax {
			flags = append(flags, "barpath_excessive_jcurve")
			verdict = "warn"
		}
	}

	if dx, ok := r.Features["drift_x_cm"]; ok {
		if thr, ok2 := thresholds["drift_err"]; ok2 && math.Abs(dx) > thr {
			flags = append(flags, "barpath_drift")
			verdict = "warn"
		}
	}

	if rms, ok := r.Features["rms_x_cm"]; ok {
		if thr, ok2 := thresholds["rms_x_warn"]; ok2 && rms > thr {
			flags = append(flags, "barpath_instability")
			verdict = "warn"
		}
	}

	r.Flags = flags
	r.Verdict = verdict
	return r
}
