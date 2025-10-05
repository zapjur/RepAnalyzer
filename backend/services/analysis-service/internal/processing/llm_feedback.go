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
- Limit "key_wins" and "key_fixes" to 2–3 items each, most impactful first.

How to write better coaching cues ("fix_cue"):
- Make it prescriptive and phase-specific: WHEN + WHAT + HOW.
  • WHEN: the phase or moment (e.g., "off the floor", "out of the hole", "on descent", "at lockout").
  • WHAT: the specific action to take.
  • HOW: the mechanism or sensation to aim for (bracing, bar position, joint action).
- Keep it concise: 1–2 short sentences (max ~20–25 words total).
- If helpful, add a very short drill or constraint at the end (e.g., "Drill: 2-count pause at mid-range", "Drill: tempo 3-0-1").
- Examples:
  • "Off the floor, push the floor away and keep hips level with shoulders; engage lats to keep the bar over midfoot."
  • "Out of the hole, drive knees forward and out while bracing; keep chest up to reduce torso lean."
  • "On descent, use a 3-sec tempo and keep the bar over the lower chest; press back slightly toward the rack on ascent."

Focus points per exercise:
- Deadlift: torso/back angle (avoid excessive rounding), hip rise vs. bar speed (hips_shoot_up), bar drift from midfoot (bar over ankle), hitching/stalling counts, lockout finish.
- Squat: depth consistency (depth_ok), torso angle at bottom, bar path over midfoot (drift_x_cm, jcurve_dx_cm), smooth tempo out of the hole.
- Bench Press: eccentric control (ecc_p95_vy_m_s), sticking/stall, bar path J-curve size (jcurve_dx_cm), lateral drift (drift_x_cm), path stability (rms_x_cm).

Issue mapping (use when flags/metrics suggest them):
- deadlift: hips_shoot_up, hitching/stall, barpath_drift, bar_not_over_midfoot, shoulders_too_far_over_bar
- squat: depth_insufficient, torso_lean_high, barpath_drift
- bench: eccentric_too_fast, stall, barpath_too_linear, barpath_excessive_jcurve, barpath_instability, barpath_drift

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
          "code": "string",    // e.g. torso_lean_high, barpath_drift, depth_insufficient, knees_caving, stall
          "evidence": "string",// include metric + threshold where relevant, e.g. "torso_angle_bottom_deg=58° > warn=55°"
          "fix_cue": "string"  // WHEN + WHAT + HOW (+ optional short Drill), e.g. "Out of the hole, drive knees forward and out while bracing; keep chest up. Drill: 2-count pause at parallel."
        }
      ]
    }
  ]
}`

	params := openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
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
