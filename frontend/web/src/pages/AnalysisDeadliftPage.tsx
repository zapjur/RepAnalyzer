import React, { useEffect, useMemo, useRef, useState } from "react";
import { useAnalyses, VideoAnalysis } from "../contexts/AnalysesContext";
import { useNavigate, useParams, Link } from "react-router-dom";
import VelocityChart from "../components/VelocityChart";

const cx = (...classes: Array<string | false | null | undefined>) => classes.filter(Boolean).join(" ");

type TabKey = "barpath" | "pose" | "analysis";

const AnalysisDeadliftPage: React.FC = () => {
    const { fetchByVideo } = useAnalyses();
    const { videoId } = useParams<{ videoId: string }>();

    const [items, setItems] = useState<VideoAnalysis[]>([]);
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
                const data = await fetchByVideo(parsedVideoId);
                if (!mounted) return;
                setItems((data || []).slice().sort((a, b) => b.id - a.id));
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

    const currentList =
        active === "barpath" ? barpathMerged : active === "pose" ? byType.pose : byType.analysis;

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
                        <Tab label="Technique Analysis (soon)" active={active === "analysis"} onClick={() => setActive("analysis")} />
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

                {!loading && !error && (Array.isArray(currentList) ? currentList.length === 0 : currentList.length === 0) ? (
                    <p>No videos yet</p>
                ) : !loading && !error ? (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        {(currentList as Array<VideoAnalysis & { mergedCsvUrl?: string | null }>).map((a) => (
                            <BarpathTile key={a.id} a={a} showChart={active === "barpath"} />
                        ))}
                    </div>
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

export default AnalysisDeadliftPage;
