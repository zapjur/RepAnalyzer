package processing

import (
	"analysis-service/internal/client"
	"analysis-service/internal/minio"
	db "analysis-service/proto"
	"analysis-service/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func px2cm(px, mpp float64) float64 { return px * mpp * 100.0 }

func nearestPoseByTime(pose []PoseSample, t, fps float64) (PoseSample, bool) {
	if len(pose) == 0 || fps <= 0 {
		return PoseSample{}, false
	}
	tol := 0.5 / fps
	best := -1
	bestDT := 1e9
	for i := range pose {
		dt := math.Abs(pose[i].T - t)
		if dt < bestDT {
			bestDT = dt
			best = i
		}
	}
	if best >= 0 && bestDT <= tol {
		return pose[best], true
	}
	return PoseSample{}, false
}

func GenerateAnalysis(ctx context.Context, minioClient *minio.Client, grpcClient *client.Client, req types.AnalysisRequest) error {
	tmpDir, err := downloadCSVToTmp(ctx, minioClient, req.Bucket, req.ObjectKey, req.Auth0Id)
	if err != nil {
		return err
	}
	posePath := filepath.Join(tmpDir, "pose.csv")
	barPath := filepath.Join(tmpDir, "barpath.csv")
	poseMetaPath := filepath.Join(tmpDir, "pose_meta.json")
	barMetaPath := filepath.Join(tmpDir, "barpath_meta.json")
	defer func() {
		_ = os.Remove(posePath)
		_ = os.Remove(barPath)
		_ = os.Remove(poseMetaPath)
		_ = os.Remove(barMetaPath)
		_ = os.RemoveAll(tmpDir)
	}()
	log.Println("Downloaded CSVs to", tmpDir)

	bar, err := readBarpathCSV(barPath)
	if err != nil {
		return fmt.Errorf("read barpath: %w", err)
	}
	ex := strings.ToLower(req.ExerciseName)
	usePose := (ex == "deadlift" || ex == "squat")

	var pose []PoseSample
	if usePose {
		pose, err = readPoseCSV(posePath)
		if err != nil {
			return fmt.Errorf("read pose: %w", err)
		}
	}

	var bm barMeta
	var pm poseMeta
	_ = readJSON(barMetaPath, &bm)
	_ = readJSON(poseMetaPath, &pm)

	var fps float64
	switch {
	case bm.FPS > 0:
		fps = bm.FPS
	case pm.FPS > 0:
		fps = pm.FPS
	default:
		fps = 25
	}

	for i := range bar {
		if bar[i].MPP <= 0 {
			if bm.MetersPerPixel > 0 {
				bar[i].MPP = bm.MetersPerPixel
			} else {
				bar[i].MPP = 0.005
			}
		}
		if bar[i].T == 0 {
			bar[i].T = float64(bar[i].Frame) / fps
		}
		if i > 0 && math.Abs(bar[i].Vy)+math.Abs(bar[i].VyS) < 1e-9 {
			dy := bar[i].Y - bar[i-1].Y
			dt := bar[i].T - bar[i-1].T
			if dt > 0 {
				vy := -(dy / dt) * bar[i].MPP
				bar[i].Vy = vy
				bar[i].VyS = vy
			}
		}
	}

	if usePose {
		for i := range pose {
			if !(pose[i].T > 0) {
				pose[i].T = float64(pose[i].Frame) / fps
			}
		}
		pose = cleanPose(pose)
	}

	bar = cleanBarpath(bar)
	reps := segmentConcentric(bar)
	if len(reps) == 0 {
		return fmt.Errorf("no reps found")
	}
	reps = mergeReps(bar, reps, fps, thresholds["seg_min_descent_cm"], thresholds["merge_gap_s"])

	for i, rc := range reps {
		dur := bar[rc.Top].T - bar[rc.Bottom].T
		dy := (bar[rc.Bottom].Y - bar[rc.Top].Y) * bar[rc.Bottom].MPP * 100.0
		log.Printf("[rep %d] dur=%.2fs, rise=%.1f cm", i+1, dur, dy)
	}

	var anchor ankleAnchor
	if usePose {
		anchor = chooseStableAnkle(bar, pose)
	}

	footLenCm := thresholds["foot_len_cm"]
	if footLenCm <= 0 {
		footLenCm = 26
	}

	out := make([]RepReport, 0, len(reps))
	for i, rc := range reps {
		r := RepReport{
			Index: i + 1,
			Features: map[string]float64{
				"rms_x_cm":     rmsLateral(bar, rc.Bottom, rc.Top),
				"drift_x_cm":   driftX(bar, rc.Bottom, rc.Top),
				"jcurve_dx_cm": jCurveDX(bar, rc.Bottom, rc.Top),
				"stall_count":  float64(stallCount(bar, rc.Bottom, rc.Top, thresholds["stall_eps"], thresholds["stall_min_s"])),
			},
		}
		if rc.PrevTop >= 0 && rc.PrevTop < rc.Bottom {
			r.Features["ecc_p95_vy_m_s"] = eccentricP95Vy(bar, rc.PrevTop, rc.Bottom)
		} else {
			r.Features["ecc_p95_vy_m_s"] = 0
		}

		if usePose && len(pose) > 0 && anchor.Samples > 0 {
			ankleX := anchor.AnkleXMedian

			tStart := bar[rc.Bottom].T
			tLock := bar[rc.Top].T
			dur := tLock - tStart

			win := thresholds["hips_window_s"]
			if win <= 0 {
				win = 0.5
			}
			tTorso := tStart + math.Min(win, dur*0.5)
			if pTorso, ok := nearestPoseByTime(pose, tTorso, fps); ok {
				r.Features["torso_angle_bottom_deg"] = torsoAngleDegAt(pTorso)
			}

			if pStart, ok := nearestPoseByTime(pose, tStart, fps); ok {
				barStartX := bar[rc.Bottom].X
				shX := 0.0
				okSh := false
				if pStart.LS.Valid && pStart.RS.Valid {
					shX = 0.5 * (pStart.LS.X + pStart.RS.X)
					okSh = true
				} else if pStart.LS.Valid {
					shX = pStart.LS.X
					okSh = true
				} else if pStart.RS.Valid {
					shX = pStart.RS.X
					okSh = true
				}
				if okSh {
					sbCm := px2cm(math.Abs(shX-barStartX), bar[rc.Bottom].MPP)
					r.Features["shoulder_bar_offset_start_cm"] = sbCm
					r.Features["shoulder_bar_offset_start_fl"] = sbCm / footLenCm
				}
			}

			barStartX := bar[rc.Bottom].X
			barLockX := bar[rc.Top].X

			axStartCm := px2cm(barStartX-ankleX, bar[rc.Bottom].MPP)
			axLockCm := px2cm(barLockX-ankleX, bar[rc.Top].MPP)

			r.Features["bar_over_ankle_start_cm"] = axStartCm
			r.Features["bar_over_ankle_lock_cm"] = axLockCm
			r.Features["bar_over_ankle_start_fl"] = axStartCm / footLenCm
			r.Features["bar_over_ankle_lock_fl"] = axLockCm / footLenCm

			if ex == "squat" {
				if pBottom, ok := nearestPoseByTime(pose, tStart, fps); ok {
					if pBottom.LH.Valid && pBottom.RH.Valid && pBottom.LK.Valid && pBottom.RK.Valid {
						hipY := 0.5 * (pBottom.LH.Y + pBottom.RH.Y)
						kneeY := 0.5 * (pBottom.LK.Y + pBottom.RK.Y)
						if hipY > kneeY {
							r.Features["depth_ok"] = 1
						} else {
							r.Features["depth_ok"] = 0
						}
					}
				}
			}

			if ex == "deadlift" {
				r.Features["hips_shoot_up_ratio"] = hipsShootUpRatio(bar, pose, rc, 0.5)
			}
		}

		switch ex {
		case "deadlift":
			r = judgeDeadlift(r)
		case "squat":
			r = judgeSquat(r)
		case "bench", "bench press", "benchpress":
			r = judgeBench(r)
		default:
			r = judgeDeadlift(r)
		}
		out = append(out, r)
	}

	rep := Report{
		VideoID:     req.VideoID,
		Exercise:    ex,
		Summary:     summarize(out),
		Reps:        out,
		LLMFeedback: nil,
		Version:     "heuristics-1.2",
		Thresholds:  thresholds,
		Meta:        map[string]any{"fps": fps},
		CreatedAt:   time.Now().UTC(),
	}

	fb, err := GenerateLLMFeedback(ctx, rep)
	if err != nil {
		log.Printf("LLM feedback error: %v", err)
	} else {
		rep.LLMFeedback = json.RawMessage(fb)
	}

	b, _ := json.Marshal(rep)
	log.Printf("ANALYSIS REPORT\n%s\n", string(b))

	vidID, err := strconv.ParseInt(req.VideoID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid video ID: %w", err)
	}

	resp, err := grpcClient.DBService.SaveAnalysisJSON(ctx, &db.SaveAnalysisJSONRequest{
		VideoId:     vidID,
		PayloadJson: string(b),
	})
	if err != nil {
		log.Printf("SaveAnalysisJSON error: %v", err)
		return fmt.Errorf("save analysis via gRPC: %w", err)
	}
	if !resp.Success {
		log.Printf("SaveAnalysisJSON failed: %s", resp.Message)
		return fmt.Errorf("save analysis failed: %s", resp.Message)
	}

	return nil
}
