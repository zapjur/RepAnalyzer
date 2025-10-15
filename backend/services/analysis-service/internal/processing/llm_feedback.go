package processing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

const cueLibraryJSON = `{
  "cue_library": {
    "_notes": {
      "purpose": "Map metrics/flags to issues, evidence text, short coaching cues, and concrete drills. Used as read-only context.",
      "placeholders": "Use {{metric}} and {{threshold}} to inject numbers, e.g., 'rms_x_cm={{metric}} cm > warn={{threshold}} cm'. Round per your rules.",
      "algo_thresholds": ["seg_vy_start","seg_vy_stop","seg_min_dur_s","seg_min_vert_cm","seg_min_descent_cm","seg_quiet_s","merge_gap_s","midfoot_calibrate","hips_window_s"],
      "do_not_surface": "Do not mention algo_thresholds in output. They are only for segmentation."
    },

    "definitions": {
      "barpath_instability": "Lateral wobble of the bar path (high RMS in X).",
      "barpath_drift": "Net horizontal shift of bar from start to finish.",
      "hips_shoot_up": "Hips rise faster than shoulders during the initial pull; torso angle increases early.",
      "torso_rounding": "Excessive spinal/torso flexion angle at bottom or through the pull.",
      "eccentric_too_fast": "Descent velocity above the control threshold (exercise-specific)."
    },

    "shared": {
      "barpath_instability": {
        "trigger": "rms_x_cm > rms_x_warn (warn) or > rms_x_err (error)",
        "evidence_template": "rms_x_cm={{metric}} cm > {{threshold}} cm",
        "fix_cues": [
          "From setup, squeeze lats and pull the bar toward you; keep forearms vertical as you move.",
          "On ascent, drive straight up; think 'bar glued to shirt.' Keep eyes fixed on one point."
        ],
        "drills": [
          "Tempo reps 3-0-1 with video from front.",
          "Isometric lat sweep: 2-second pre-pull against the bar.",
          "Pause halfway up (1–2 s) to hold vertical line."
        ]
      },
      "barpath_drift": {
        "trigger": "abs(drift_x_cm) > drift_err",
        "evidence_template": "drift_x_cm={{metric}} cm > err={{threshold}} cm",
        "fix_cues": [
          "Through the rep, stack bar over midfoot; press up/down, not forward/back.",
          "On the way up, push evenly through the whole foot; keep big toe and heel planted."
        ],
        "drills": [
          "Tempo 3-0-1 focusing on balance over midfoot.",
          "Tape a vertical 'lane' on the floor and stay inside it.",
          "Wall cue: pull 5–8 cm in front of a wall to prevent drift."
        ]
      },
      "barpath_too_linear": {
        "trigger": "jcurve_dx_cm < jcurve_min (bench-specific)",
        "evidence_template": "jcurve_dx_cm={{metric}} cm < min={{threshold}} cm",
        "fix_cues": [
          "Off the chest, aim slightly back toward the rack, not straight up; keep elbows under the bar."
        ],
        "drills": [
          "Pin press at mid-range to groove back-then-up path.",
          "Spoto press (1 s pause 1–2 cm above chest)."
        ]
      },
      "barpath_excessive_jcurve": {
        "trigger": "jcurve_dx_cm > jcurve_max",
        "evidence_template": "jcurve_dx_cm={{metric}} cm > max={{threshold}} cm",
        "fix_cues": [
          "Out of the bottom, keep the path tighter; push up and slightly back, not in an arc."
        ],
        "drills": [
          "Board press / pin press to shorten the arc.",
          "Light technique sets focusing on straight-up after initial back drive."
        ]
      },

    "deadlift": {
      "bar_not_over_midfoot": {
        "trigger": "midfoot_offset_cm > midfoot_offset_cm_deadlift",
        "evidence_template": "midfoot_offset_cm={{metric}} cm > {{threshold}} cm",
        "fix_cues": [
          "At setup, bring shins to the bar and pull slack; bar over midfoot before you break the floor.",
          "Off the floor, push the floor away and keep shoulders just in front of the bar."
        ],
        "drills": [
          "Floating deadlift (bar hovers 1–2 cm before pull) to lock midfoot balance.",
          "Block pulls from mid-shin to practice vertical drive."
        ]
      },
      "shoulders_too_far_over_bar": {
        "trigger": "shoulder_over_bar_fraction > shoulder_over_bar_fl",
        "evidence_template": "shoulder_over_bar_fraction={{metric}} > {{threshold}}",
        "fix_cues": [
          "At setup, sit slightly back and lock lats; balance over midfoot so shoulders aren’t too far ahead.",
          "Keep chest tall and drag the bar up the legs."
        ],
        "drills": [
          "Lat prep: straight-arm pulldown 2×10 before sets.",
          "Touch-the-legs cue: graze shins/thighs the whole way up."
        ]
      },
      "hips_shoot_up": {
        "trigger": "hips_shoot_up flag present OR clear early hip rise within hips_window_s while bar vy is low",
        "evidence_template": "early hip rise (hips↑ vs shoulders within {{threshold}} s)",
        "fix_cues": [
          "Push the floor away with your legs; keep hips and shoulders rising together. Think 'legs drive, bar moves as one.'",
   	  	  "Stay wedged and tight before the pull; keep the same back angle until the bar leaves the floor."
        ],
        "drills": [
          "Paused deadlift 2–3 cm off the floor (1–2 s).",
          "Tempo 2-0-2 from floor to knee focusing on same torso angle.",
          "Leg press cueing on full-foot pressure between sets."
        ]
      },
      "torso_rounding": {
        "trigger": "torso_angle_bottom_deg > torso_dl_err",
        "evidence_template": "torso_angle_bottom_deg={{metric}}° > err={{threshold}}°",
        "fix_cues": [
          "Before pull, brace 360° and lock lats; keep chest up as the bar breaks the floor.",
          "Show your logo forward as you push; keep the bar close to shins."
        ],
        "drills": [
          "McGill big-3 between warm-ups.",
          "Tempo 3-0-1 off the floor with video side view.",
          "Belted cueing: expand into the belt before every rep."
        ]
      },
      "eccentric_too_fast": {
        "trigger": "ecc_p95_vy_m_s > ecc_vy_warn_dl",
        "evidence_template": "ecc_p95_vy_m_s={{metric}} m/s > warn={{threshold}} m/s",
        "fix_cues": [
          "On the way down, control the bar to the floor; keep it on your legs.",
          "Guide the bar down your thighs; don’t drop it."
        ],
        "drills": [
          "2-second controlled negatives.",
          "Touch-and-go sets with strict thigh contact."
        ]
      }
    },

    "squat": {
      "depth_insufficient": {
        "trigger": "depth_ok = false (your detector); do not use for deadlift.",
        "evidence_template": "depth flag failed on rep (depth_ok=false)",
        "fix_cues": [
          "On descent, sit straight down between the hips while keeping chest up; stay tight and hit depth.",
          "Use a slightly slower 3-sec descent to find consistent depth."
        ],
        "drills": [
          "1–2 s pause in the hole.",
          "Box squat to just-below-parallel height.",
          "Tempo 3-0-1 with depth marker."
        ]
      },
      "torso_lean_high": {
        "trigger": "torso_angle_bottom_deg > torso_sq_warn (warn) or > torso_sq_err (error)",
        "evidence_template": "torso_angle_bottom_deg={{metric}}° > {{threshold}}°",
        "fix_cues": [
          "Out of the hole, drive knees forward and out while bracing; keep chest up.",
          "Keep elbows under the bar and ribcage down; stand straight up."
        ],
        "drills": [
          "High-bar pause squats 2 s.",
          "Counterbalance goblet squats to learn upright torso.",
          "Tempo 3-0-1 with front-facing video."
        ]
      },
      "bar_not_over_midfoot": {
        "trigger": "midfoot_offset_cm > midfoot_offset_cm_squat",
        "evidence_template": "midfoot_offset_cm={{metric}} cm > {{threshold}} cm",
        "fix_cues": [
          "During descent, track knees over toes and keep the bar over midfoot.",
          "On ascent, push evenly through the whole foot; avoid rocking to toes/heels."
        ],
        "drills": [
          "Heels-to-wall squats (light) to prevent forward drift.",
          "Tempo 3-0-1 focusing on balance tripod foot."
        ]
      },
      "barpath_drift_or_instability": {
        "trigger": "Use shared.barpath_instability or shared.barpath_drift",
        "evidence_template": "see shared templates",
        "fix_cues": [
          "Brace and drive straight up; keep gaze fixed.",
          "Squeeze upper back; pin the bar to a vertical track."
        ],
        "drills": [
          "Tempo squats 3-0-1 with bar path overlay if available.",
          "Pin squats (1–2 s) just below sticking point."
        ]
      },
      "eccentric_too_fast": {
        "trigger": "ecc_p95_vy_m_s > ecc_vy_warn_sq",
        "evidence_template": "ecc_p95_vy_m_s={{metric}} m/s > warn={{threshold}} m/s",
        "fix_cues": [
          "On descent, use a 3-sec tempo and keep tension; arrive balanced over midfoot.",
          "Control the last third of the drop; keep hips under the bar."
        ],
        "drills": [
          "3-0-1 tempo sets with metronome.",
          "Pause squats 1–2 s in the hole."
        ]
      }
    },

    "bench": {
      "eccentric_too_fast": {
        "trigger": "ecc_p95_vy_m_s > ecc_vy_warn_bench",
        "evidence_template": "ecc_p95_vy_m_s={{metric}} m/s > warn={{threshold}} m/s",
        "fix_cues": [
          "On descent, lower with control to the lower chest; keep shoulder blades squeezed.",
          "Touch softly, keep tension, then press."
        ],
        "drills": [
          "3-0-1 tempo benches.",
          "Spoto press (1 s pause above chest)."
        ]
      },
      "barpath_too_linear": {
        "trigger": "jcurve_dx_cm < jcurve_min",
        "evidence_template": "jcurve_dx_cm={{metric}} cm < min={{threshold}} cm",
        "fix_cues": [
          "From the chest, press slightly back (toward eyes) then up; follow a subtle J path."
        ],
        "drills": [
          "Board/pin press to learn back-then-up groove.",
          "Light technique sets focusing on elbow under bar."
        ]
      },
      "barpath_excessive_jcurve": {
        "trigger": "jcurve_dx_cm > jcurve_max",
        "evidence_template": "jcurve_dx_cm={{metric}} cm > max={{threshold}} cm",
        "fix_cues": [
          "Keep the path tighter; press more vertically after the initial back drive."
        ],
        "drills": [
          "Pin press from mid-range to reduce arc.",
          "Touch-point consistency drill with video lines."
        ]
      },
      "barpath_drift_or_instability": {
        "trigger": "Use shared.barpath_instability or shared.barpath_drift",
        "evidence_template": "see shared templates",
        "fix_cues": [
          "Keep wrists stacked and forearms vertical; pull the bar out of the rack to set the groove.",
          "Lock scapulae down/back; use a consistent lower-chest touch point."
        ],
        "drills": [
          "Longer unrack hold (2 s) to set lats and groove.",
          "Feet anchored / leg drive practice sets."
        ]
      }
    },

    "calibration_and_context": {
      "foot_length_reference": {
        "trigger": "foot_len_cm used to estimate midfoot and stance (contextual)",
        "evidence_template": "foot_len_cm={{metric}} cm (context for midfoot thresholds)",
        "fix_cues": [
          "Set stance so bar starts over midfoot (about mid-arch). Use your shoe length as a reference."
        ],
        "drills": [
          "Balance drill: tripod foot (big toe, little toe, heel) under light load.",
          "Slow eccentrics while watching bar-over-midfoot on video."
        ]
      }
    }
  }
}`

func GenerateLLMFeedback(ctx context.Context, rep Report) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	inputJSON, err := json.Marshal(rep)
	if err != nil {
		return "", fmt.Errorf("marshal report: %w", err)
	}

	systemPrompt := `You are a professional strength coach.
You will receive structured analysis data (per-rep metrics and flags) for barbell lifts.
Your job is to interpret these numbers and generate clear, actionable feedback for the athlete.

General rules:
- Return ONLY valid JSON in the schema below (no prose, no extra keys).
- Write in simple, direct English, as if coaching the lifter.
- Always highlight BOTH positives ("key_wins") and fixes ("key_fixes").
- In "evidence", cite the actual metric value and, when relevant, the threshold used.
- Do NOT invent issues — base everything strictly on the provided data and flags.
- Tailor feedback to the exercise type (deadlift, squat, bench press).
- Keep numbers readable: round cm to 1 decimal, m/s to 2 decimals, angles to 0–1 decimals, and include units (cm, m/s, °).
- Limit "key_wins" and "key_fixes" to 3–4 items each, most impactful first.

How to write better coaching cues ("fix_cue"):
- Make it prescriptive and phase-specific: WHEN + WHAT + HOW.
  • WHEN: the phase or moment (e.g., "off the floor", "out of the hole", "on descent", "at lockout").
  • WHAT: the specific action to take.
  • HOW: the mechanism or sensation to aim for (bracing, bar position, joint action).
- Keep it concise: 2–3 short sentences (max ~30–40 words total).
- If helpful, add a very short drill to fix_cue or constraint at the end (e.g., "Drill: 2-count pause at mid-range", "Drill: tempo 3-0-1").
- Examples:
  • "Off the floor, push the floor away and keep hips level with shoulders; engage lats to keep the bar over midfoot."
  • "Out of the hole, drive knees forward and out while bracing; keep chest up to reduce torso lean."
  • "On descent, use a 3-sec tempo and keep the bar over the lower chest; press back slightly toward the rack on ascent."

Issue mapping:
- deadlift: hips_shoot_up, barpath_drift, barpath_instability, bar_not_over_midfoot, shoulders_too_far_over_bar, torso_lean_high, eccentric_too_fast
- squat: depth_insufficient, torso_lean_high, bar_not_over_midfoot, barpath_drift, barpath_instability, eccentric_too_fast
- bench: eccentric_too_fast, barpath_too_linear, barpath_excessive_jcurve, barpath_instability, barpath_drift

Grading guidance (overall.grade):
- "ok": no critical flags; minor or zero warnings.
- "warn": one or more meaningful warnings; technique requires attention.
- "error": multiple severe issues or the same major issue repeated across reps.

Output schema (must follow exactly):
{
  "video_id": "string",
  "exercise": "string",
  "overall": {
    "grade": "ok|warn|error",
    "one_line_summary": "string",
    "key_wins": ["string"],
    "key_fixes": ["string"]
  },
  "rep_feedback": [
    {
      "rep": "int",
      "verdict": "ok|warn|error",
      "issues": [
        {
          "code": "string",
          "evidence": "string",
          "fix_cue": "string"
        }
      ]
    }
  ]
}

---
Inputs You’ll Receive:
1) cue_library: JSON object with triggers, evidence templates, and fix cues (read-only).
2) report: JSON object with video_id, exercise, and per-rep features/flags (provided by the user message).

What To Do:
- Use cue_library triggers to decide which issues apply, based ONLY on the report.
- Use evidence_template to build “evidence” with actual values and the relevant threshold (warn/error). Round per rules.
- For each chosen issue, select 1 best fix_cue (max 2 if strongly helpful). Keep cues phase-specific (WHEN+WHAT+HOW).
- Build “overall” from the most common/impactful issues and wins across reps.
- IMPORTANT: Do NOT include cue_library in your output. Return ONLY the JSON that matches the schema.

Resolver Rules:
- If a trigger has warn and error levels, use the highest level satisfied and cite that threshold in evidence.
- Prefer the first fix_cue unless another better matches the rep phase.
- If no triggers fire for a rep, return "issues": [] for that rep.
- Output must be valid JSON matching the schema exactly. No extra keys or prose.

Exercise-specific exclusions (hard rules):
- Depth is ONLY relevant for squat. For deadlift and bench, IGNORE any depth-related metrics/flags and DO NOT mention depth in key_wins, key_fixes, or issues.
- If exercise != "squat", you must not output any text containing the word "depth".
- J-curve issues only for bench.`

	params := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.SystemMessage(cueLibraryJSON),
			openai.UserMessage(string(inputJSON)),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		},
		Temperature: openai.Float(0.2),
		MaxTokens:   openai.Int(800),
	}

	resp, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("llm call failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	raw := resp.Choices[0].Message.Content

	var parsed any
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return "", fmt.Errorf("model did not return valid JSON: %w; raw: %s", err, raw)
	}
	pretty, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("re-marshal json: %w", err)
	}

	log.Printf("LLM feedback generated (English)\n%s\n", string(pretty))
	return string(pretty), nil
}
