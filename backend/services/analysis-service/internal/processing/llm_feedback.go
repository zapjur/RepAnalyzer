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
- Return ONLY valid JSON in the schema below.
- Use short, simple English coaching cues, as if speaking to the lifter.
- Always highlight BOTH positives ("key_wins") and fixes ("key_fixes").
- In "evidence", cite the actual metric value and the threshold if relevant.
- Do NOT invent issues — base everything strictly on the provided data and flags.
- Tailor your feedback to the exercise type (deadlift, squat, or bench press).

Focus points per exercise:
- **Deadlift**: back angle (avoid excessive rounding), hip rise vs. bar speed, bar drift away from midfoot (look for bar over ankle), lockout quality, hitching/stalling.
- **Squat**: depth consistency (hips below knees), torso lean at the bottom, bar path over midfoot, smooth tempo out of the hole.
- **Bench Press**: bar path curve (J-curve), touch point consistency, elbow flare, range of motion, bar speed (avoid stalling halfway).

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
          "code": "string",    // e.g. torso_lean_high, barpath_drift, depth_shallow, knees_caving, stall
          "evidence": "string",// include metric + threshold, e.g. "torso_angle_bottom_deg=58° > warn=55°"
          "fix_cue": "string"  // short coaching cue, e.g. "Keep chest up and brace core"
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
