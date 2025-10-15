package processing

import (
	"encoding/json"
	"time"
)

type Point struct {
	X, Y  float64
	Valid bool
}

type BarSample struct {
	Frame int
	T     float64
	X     float64
	Y     float64
	Vy    float64
	VyS   float64
	MPP   float64
}

type PoseSample struct {
	Frame int
	T     float64
	LS    Point
	RS    Point
	LH    Point
	RH    Point
	LK    Point
	RK    Point
	LA    Point
	RA    Point
}

type RepCut struct {
	Bottom  int
	Top     int
	PrevTop int
}

type RepReport struct {
	Index    int                `json:"index"`
	Verdict  string             `json:"verdict,omitempty"`
	Flags    []string           `json:"flags,omitempty"`
	Features map[string]float64 `json:"features"`
}

type Report struct {
	VideoID     string             `json:"video_id"`
	Exercise    string             `json:"exercise"`
	Summary     map[string]int     `json:"summary"`
	Reps        []RepReport        `json:"reps"`
	LLMFeedback json.RawMessage    `json:"llm_feedback,omitempty"`
	Version     string             `json:"version"`
	Thresholds  map[string]float64 `json:"thresholds"`
	Meta        map[string]any     `json:"meta"`
	CreatedAt   time.Time          `json:"created_at"`
}

var thresholds = map[string]float64{
	"seg_vy_start":       0.06,
	"seg_vy_stop":        0.02,
	"seg_min_dur_s":      0.35,
	"seg_min_vert_cm":    12,
	"seg_min_descent_cm": 20,
	"seg_quiet_s":        0.12,
	"merge_gap_s":        0.22,

	"rms_x_warn": 3,
	"rms_x_err":  6,
	"drift_err":  4,
	"jcurve_min": 1,
	"jcurve_max": 8,

	"torso_dl_err":  60,
	"torso_sq_warn": 55,
	"torso_sq_err":  60,

	"ecc_vy_warn_dl":    1.0,
	"ecc_vy_warn_sq":    0.8,
	"ecc_vy_warn_bench": 0.6,

	"foot_len_cm":                26,
	"midfoot_offset_cm_deadlift": 6,
	"midfoot_offset_cm_squat":    5,
	"midfoot_offset_cm_default":  6,
	"shoulder_over_bar_fl":       0.25,

	"midfoot_calibrate": 0,

	"hips_window_s": 0.5,
}

func summarize(reps []RepReport) map[string]int {
	out := map[string]int{"ok": 0, "warn": 0, "error": 0}
	for _, r := range reps {
		switch r.Verdict {
		case "error":
			out["error"]++
		case "warn":
			out["warn"]++
		default:
			out["ok"]++
		}
	}
	return out
}
