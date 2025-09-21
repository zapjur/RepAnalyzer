package processing

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func readJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

type barMeta struct {
	FPS            float64                      `json:"fps"`
	Width          int                          `json:"width"`
	Height         int                          `json:"height"`
	Origin         struct{ X, Y, Units string } `json:"origin"`
	MetersPerPixel float64                      `json:"meters_per_pixel"`
	Notes          string                       `json:"notes"`
}

type poseMeta struct {
	FPS         float64                      `json:"fps"`
	Width       int                          `json:"width"`
	Height      int                          `json:"height"`
	Origin      struct{ X, Y, Units string } `json:"origin"`
	Dataset     string                       `json:"dataset"`
	Indices     map[string]int               `json:"indices"`
	BodyWritten map[string]int               `json:"body_indices_written"`
	Notes       string                       `json:"notes"`
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

func readBarpathCSV(path string) ([]BarSample, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	hdr, err := r.Read()
	if err != nil {
		return nil, err
	}

	idx := map[string]int{}
	for i, h := range hdr {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}
	get := func(key string) (int, bool) {
		i, ok := idx[strings.ToLower(key)]
		return i, ok
	}

	var out []BarSample
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) == 0 {
			continue
		}

		var s BarSample
		if i, ok := get("frame"); ok {
			s.Frame = int(parseFloat(rec[i]))
		}
		if i, ok := get("t"); ok {
			s.T = parseFloat(rec[i])
		}
		if i, ok := get("x_px"); ok {
			s.X = parseFloat(rec[i])
		}
		if i, ok := get("y_px"); ok {
			s.Y = parseFloat(rec[i])
		}
		if i, ok := get("vy_m_s"); ok {
			s.Vy = parseFloat(rec[i])
		}
		if i, ok := get("vy_smooth_m_s"); ok {
			s.VyS = parseFloat(rec[i])
		}
		if i, ok := get("meters_per_pixel"); ok {
			s.MPP = parseFloat(rec[i])
		}

		out = append(out, s)
	}
	return out, nil
}

func readPoseCSV(path string) ([]PoseSample, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	hdr, err := r.Read()
	if err != nil {
		return nil, err
	}

	findCol := func(cands ...string) int {
		for _, want := range cands {
			for i, h := range hdr {
				if strings.EqualFold(strings.TrimSpace(h), want) {
					return i
				}
			}
		}
		return -1
	}

	fi := findCol("frame", "frame_idx")
	ti := findCol("t", "time", "timestamp")
	if fi < 0 {
		return nil, fmt.Errorf("pose csv missing frame/frame_idx")
	}

	type xy struct{ X, Y int }
	col := func(name string) xy {
		xName := name + "_x"
		yName := name + "_y"
		xi, yi := -1, -1
		for i, h := range hdr {
			hc := strings.TrimSpace(h)
			if strings.EqualFold(hc, xName) {
				xi = i
			} else if strings.EqualFold(hc, yName) {
				yi = i
			}
		}
		return xy{X: xi, Y: yi}
	}

	cols := map[string]xy{
		"Left Shoulder":  col("Left Shoulder"),
		"Right Shoulder": col("Right Shoulder"),
		"Left Hip":       col("Left Hip"),
		"Right Hip":      col("Right Hip"),
		"Left Knee":      col("Left Knee"),
		"Right Knee":     col("Right Knee"),
		"Left Ankle":     col("Left Ankle"),
		"Right Ankle":    col("Right Ankle"),
	}

	readPt := func(rec []string, xi, yi int) Point {
		if xi < 0 || yi < 0 || xi >= len(rec) || yi >= len(rec) {
			return Point{}
		}
		x := parseFloat(rec[xi])
		y := parseFloat(rec[yi])
		if x == 0 && y == 0 {
			return Point{}
		}
		return Point{X: x, Y: y, Valid: true}
	}

	var out []PoseSample
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		ps := PoseSample{
			Frame: int(parseFloat(rec[fi])),
		}
		if ti >= 0 {
			ps.T = parseFloat(rec[ti])
		} else {
			ps.T = -1
		}

		ps.LS = readPt(rec, cols["Left Shoulder"].X, cols["Left Shoulder"].Y)
		ps.RS = readPt(rec, cols["Right Shoulder"].X, cols["Right Shoulder"].Y)
		ps.LH = readPt(rec, cols["Left Hip"].X, cols["Left Hip"].Y)
		ps.RH = readPt(rec, cols["Right Hip"].X, cols["Right Hip"].Y)
		ps.LK = readPt(rec, cols["Left Knee"].X, cols["Left Knee"].Y)
		ps.RK = readPt(rec, cols["Right Knee"].X, cols["Right Knee"].Y)
		ps.LA = readPt(rec, cols["Left Ankle"].X, cols["Left Ankle"].Y)
		ps.RA = readPt(rec, cols["Right Ankle"].X, cols["Right Ankle"].Y)

		out = append(out, ps)
	}
	return out, nil
}
