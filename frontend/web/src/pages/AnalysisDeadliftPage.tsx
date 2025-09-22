import React, { useEffect, useMemo, useRef, useState } from "react";
import { useAnalyses, VideoAnalysis } from "../contexts/AnalysesContext";
import { useNavigate, useParams, Link } from "react-router-dom";
import VelocityChart from "../components/VelocityChart";

const cx = (...classes: Array<string | false | null | undefined>) => classes.filter(Boolean).join(" ");

type TabKey = "barpath" | "pose" | "analysis";

/** ---- Types for the analysis JSON (approx) ---- */
type LLMIssue = { code?: string; evidence?: string; fix_cue?: string };
type LLMRepFeedback = { rep?: number | string; verdict?: string; issues?: LLMIssue[] };
type LLMOverall = { grade?: string; one_line_summary?: string; key_wins?: string[]; key_fixes?: string[] };
type LLMFeedback = { overall?: LLMOverall; rep_feedback?: LLMRepFeedback[]; exercise?: string; video_id?: string };
type AnalysisRep = { index: number; verdict: string; flags?: string[]; features?: Record<string, number> };
type HeuristicsSummary = { ok?: number; warn?: number; error?: number };
type AnalysisPayload = {
    meta?: { fps?: number };
    reps?: AnalysisRep[];
    summary?: HeuristicsSummary;
    version?: string;
    exercise?: string;
    video_id?: string;
    created_at?: string;
    thresholds?: Record<string, number>;
    llm_feedback?: LLMFeedback;
};

const AnalysisDeadliftPage: React.FC = () => {
    const { fetchByVideo } = useAnalyses();
    const { videoId } = useParams<{ videoId: string }>();

    const [items, setItems] = useState<VideoAnalysis[]>([]);
    const [analysis, setAnalysis] = useState<AnalysisPayload | null>(null);

    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [active, setActive] = useState<TabKey>("barpath");

    const navigate = useNavigate();
    const parsedVideoId = Number(videoId);

    useEffect(() => {
        let mounted = true;
        (async () => {
            if (!videoId || Number.isNaN(parsedVideoId)) {
                setError("Invalid or missing videoId in the URL.");
                setLoading(false);
                return;
            }
            try {
                setLoading(true);
                setError(null);
                const { videos, analysis } = await fetchByVideo(parsedVideoId);
                if (!mounted) return;
                const list = (videos || []).slice().sort((a, b) => b.id - a.id);
                setItems(list);
                setAnalysis((analysis ?? null) as AnalysisPayload | null);
            } catch (e: any) {
                if (!mounted) return;
                setError(e?.message || "Failed to fetch analyses.");
            } finally {
                if (mounted) setLoading(false);
            }
        })();
        return () => {
            mounted = false;
        };
    }, [fetchByVideo, parsedVideoId, videoId]);

    const byType = useMemo(() => {
        const groups: Record<TabKey, VideoAnalysis[]> = { barpath: [], pose: [], analysis: [] } as const as Record<
            TabKey,
            VideoAnalysis[]
        >;
        for (const it of items) {
            const type = it.type?.toLowerCase?.() || "";
            if (type.includes("bar")) groups.barpath.push(it);
            else if (type.includes("pose")) groups.pose.push(it);
            else groups.analysis.push(it);
        }
        return groups;
    }, [items]);

    const allCsvByVideoId = useMemo(() => {
        const m = new Map<number, string>();
        for (const it of items) {
            if (it.csv_url) m.set(it.video_id, it.csv_url);
            if (it.url && /\.csv($|\?)/i.test(it.url)) m.set(it.video_id, it.url);
        }
        return m;
    }, [items]);

    const barpathMerged = useMemo(() => {
        const list = byType.barpath;
        const acc = new Map<number, { video?: VideoAnalysis; csvUrl?: string }>();
        for (const it of list) {
            const e = acc.get(it.video_id) ?? {};
            if (it.url && /\.mp4($|\?)/i.test(it.url)) {
                e.video = it;
                if (it.csv_url) e.csvUrl = it.csv_url;
            }
            if (it.csv_url) e.csvUrl = it.csv_url;
            if (it.url && /\.csv($|\?)/i.test(it.url) && !e.csvUrl) e.csvUrl = it.url;
            if (!e.csvUrl) e.csvUrl = allCsvByVideoId.get(it.video_id) ?? e.csvUrl;
            acc.set(it.video_id, e);
        }
        const merged: Array<VideoAnalysis & { mergedCsvUrl: string | null }> = [];
        for (const [, { video, csvUrl }] of acc) {
            if (!video) continue;
            merged.push({ ...video, mergedCsvUrl: csvUrl ?? null });
        }
        merged.sort((a, b) => b.id - a.id);
        return merged;
    }, [byType.barpath, allCsvByVideoId]);

    /** only used for grids; analysis tab has its own panel */
    const currentList = active === "barpath" ? barpathMerged : active === "pose" ? byType.pose : [];

    return (
        <div className="flex h-screen">
            <aside className="w-64 bg-stone-50 text-white p-6 h-screen flex flex-col">
                <Link to={"/dashboard"}>
                    <img src="/repanalyzer-logo-small.png" alt="Logo" className="w-32 mx-auto mb-4" />
                </Link>

                <div className="flex flex-col flex-1">
                    <div className="mb-6">
                        <Tab label="Bar Path" active={active === "barpath"} onClick={() => setActive("barpath")} />
                        <Tab label="Pose Estimation" active={active === "pose"} onClick={() => setActive("pose")} />
                        <Tab label="Technique Analysis" active={active === "analysis"} onClick={() => setActive("analysis")} />
                    </div>
                </div>

                <button
                    onClick={() => navigate(-1)}
                    className="bg-red-500 text-stone-950 block p-2 rounded-full hover:bg-red-600 mt-auto"
                >
                    ← Back
                </button>
            </aside>

            <main className="flex-1 p-6 bg-stone-200">
                {loading && <p>Loading…</p>}
                {error && <p className="text-red-600">Error: {error}</p>}

                {!loading && !error && active !== "analysis" && currentList.length === 0 ? (
                    <p>No videos yet</p>
                ) : !loading && !error && active !== "analysis" ? (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        {(currentList as Array<VideoAnalysis & { mergedCsvUrl?: string | null }>).map((a) => (
                            <BarpathTile key={a.id} a={a} showChart={active === "barpath"} />
                        ))}
                    </div>
                ) : null}

                {!loading && !error && active === "analysis" ? (
                    <TechniqueAnalysisPanel analysis={analysis} />
                ) : null}
            </main>
        </div>
    );
};

function Tab({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
    return (
        <button
            onClick={onClick}
            className={cx("text-stone-950 block p-2 rounded mb-2 w-full text-left", active ? "bg-stone-300" : "hover:bg-stone-300")}
        >
            {label}
        </button>
    );
}

function BarpathTile({
                         a,
                         showChart,
                     }: {
    a: VideoAnalysis & { mergedCsvUrl?: string | null };
    showChart: boolean;
}) {
    const videoRef = useRef<HTMLVideoElement | null>(null);

    return (
        <div className="flex flex-col items-center">
            <video ref={videoRef} controls className="w-full mb-2">
                <source src={a.url} type="video/mp4" />
                Your browser does not support the video tag.
            </video>

            {showChart && a.mergedCsvUrl ? (
                <div className="w-full">
                    <VelocityChart csvUrl={a.mergedCsvUrl} videoRef={videoRef} useSmoothedVelocity height={220} />
                </div>
            ) : null}
        </div>
    );
}

/** ---- Analysis tab components ---- */

function Badge({ children, tone = "neutral" }: { children: React.ReactNode; tone?: "ok" | "warn" | "error" | "neutral" }) {
    const map: Record<string, string> = {
        ok: "bg-green-200 text-green-900",
        warn: "bg-yellow-200 text-yellow-900",
        error: "bg-red-200 text-red-900",
        neutral: "bg-stone-300 text-stone-900",
    };
    return <span className={`px-2 py-1 rounded text-xs font-medium ${map[tone]}`}>{children}</span>;
}

function TechniqueAnalysisPanel({ analysis }: { analysis: AnalysisPayload | null }) {
    if (!analysis) return <p>No analysis yet</p>;

    const lf = analysis.llm_feedback;
    const overall = lf?.overall;
    const grade = (overall?.grade || "").toLowerCase() as "ok" | "warn" | "error" | "";
    const tone: "ok" | "warn" | "error" | "neutral" = grade === "ok" ? "ok" : grade === "warn" ? "warn" : grade === "error" ? "error" : "neutral";

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="bg-white rounded-xl p-4 shadow">
                <div className="flex items-center justify-between gap-4">
                    <div>
                        <h2 className="text-xl font-semibold text-stone-900">
                            {analysis.exercise?.toUpperCase?.() || "TECHNIQUE ANALYSIS"}{" "}
                            {analysis.created_at ? <span className="text-stone-500 text-sm">({new Date(analysis.created_at).toLocaleString()})</span> : null}
                        </h2>
                        <p className="text-stone-700 mt-1">{overall?.one_line_summary || "Technique review based on heuristics and LLM feedback."}</p>
                    </div>
                    <Badge tone={tone}>{overall?.grade?.toUpperCase?.() || "N/A"}</Badge>
                </div>

                {/* Key wins / fixes */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
                    <div className="bg-green-50 rounded-lg p-3">
                        <h3 className="font-medium text-green-900 mb-2">Key wins</h3>
                        <ul className="list-disc ml-5 text-green-900">
                            {(overall?.key_wins && overall.key_wins.length > 0 ? overall.key_wins : ["Consistent depth/ROM or bar control (auto-detected)."]).map((w, i) => (
                                <li key={i}>{w}</li>
                            ))}
                        </ul>
                    </div>
                    <div className="bg-yellow-50 rounded-lg p-3">
                        <h3 className="font-medium text-yellow-900 mb-2">Key fixes</h3>
                        <ul className="list-disc ml-5 text-yellow-900">
                            {(overall?.key_fixes && overall.key_fixes.length > 0 ? overall.key_fixes : ["Keep bar over midfoot", "Brace harder, control torso angle"]).map((f, i) => (
                                <li key={i}>{f}</li>
                            ))}
                        </ul>
                    </div>
                </div>

                {/* Heuristic summary */}
                <div className="mt-4 text-sm text-stone-700">
                    <span className="mr-3">Heuristic summary:</span>
                    <Badge tone="ok">OK: {analysis.summary?.ok ?? 0}</Badge>{" "}
                    <Badge tone="warn">WARN: {analysis.summary?.warn ?? 0}</Badge>{" "}
                    <Badge tone="error">ERROR: {analysis.summary?.error ?? 0}</Badge>{" "}
                    {analysis.meta?.fps ? <span className="ml-3 text-stone-500">FPS: {analysis.meta.fps}</span> : null}
                    {analysis.version ? <span className="ml-3 text-stone-500">ver: {analysis.version}</span> : null}
                </div>
            </div>

            {/* Per-rep feedback */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                {(analysis.reps || []).map((r) => {
                    const repFeedback: LLMRepFeedback[] = Array.isArray(lf?.rep_feedback) ? lf!.rep_feedback! : [];
                    const llmForRep = repFeedback.find((x) => Number(x.rep) === r.index);
                    const issues = llmForRep?.issues || [];
                    const repVerdict = (llmForRep?.verdict || r.verdict || "").toLowerCase() as "ok" | "warn" | "error" | "";
                    const repTone: "ok" | "warn" | "error" | "neutral" =
                        repVerdict === "ok" ? "ok" : repVerdict === "warn" ? "warn" : repVerdict === "error" ? "error" : "neutral";

                    return (
                        <div key={r.index} className="bg-white rounded-xl p-4 shadow">
                            <div className="flex items-center justify-between mb-2">
                                <h4 className="font-semibold text-stone-900">Rep {r.index}</h4>
                                <Badge tone={repTone}>{(llmForRep?.verdict || r.verdict || "N/A").toUpperCase()}</Badge>
                            </div>

                            {/* LLM issues (preferred) */}
                            {issues.length > 0 ? (
                                <ul className="space-y-2">
                                    {issues.map((it, i) => (
                                        <li key={i} className="text-sm">
                                            <div className="font-medium text-stone-900">{it.code}</div>
                                            {it.evidence ? <div className="text-stone-600">Evidence: {it.evidence}</div> : null}
                                            {it.fix_cue ? <div className="text-stone-800">Cue: {it.fix_cue}</div> : null}
                                        </li>
                                    ))}
                                </ul>
                            ) : (
                                /* Fallback: heuristics flags + a few key features */
                                <div className="text-sm">
                                    {r.flags && r.flags.length > 0 ? (
                                        <div className="mb-2">
                                            <div className="text-stone-900 font-medium">Flags</div>
                                            <div className="flex flex-wrap gap-2 mt-1">
                                                {r.flags.map((f) => (
                                                    <Badge key={f} tone="neutral">{f}</Badge>
                                                ))}
                                            </div>
                                        </div>
                                    ) : null}

                                    {r.features ? (
                                        <div>
                                            <div className="text-stone-900 font-medium">Key features</div>
                                            <div className="grid grid-cols-2 gap-1 mt-1 text-stone-700">
                                                {pickFeature(r.features, "jcurve_dx_cm")}
                                                {pickFeature(r.features, "drift_x_cm")}
                                                {pickFeature(r.features, "rms_x_cm")}
                                                {pickFeature(r.features, "stall_count")}
                                                {pickFeature(r.features, "torso_angle_bottom_deg")}
                                            </div>
                                        </div>
                                    ) : null}
                                </div>
                            )}
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

function pickFeature(features: Record<string, number>, key: string) {
    if (!(key in features)) return null;
    return (
        <div className="flex items-center justify-between">
            <span className="text-stone-500">{key}</span>
            <span className="font-mono">{Number(features[key]).toFixed(2)}</span>
        </div>
    );
}

export default AnalysisDeadliftPage;
