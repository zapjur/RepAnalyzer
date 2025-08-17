import React, { useEffect, useState } from "react";
import { useAnalyses, VideoAnalysis } from "../contexts/AnalysesContext";
import { useNavigate, useParams } from "react-router-dom";

export default function AnalysisDeadliftPage() {
    const { fetchByVideo } = useAnalyses();
    const { videoId } = useParams<{ videoId: string }>();
    const [items, setItems] = useState<VideoAnalysis[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
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

        return () => { mounted = false; };
    }, [fetchByVideo, parsedVideoId, videoId]);

    return (
        <div className="p-4">
            <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-bold">
                    Deadlift — Analyses (video #{!Number.isNaN(parsedVideoId) ? parsedVideoId : "?"})
                </h2>
                <button
                    onClick={() => navigate(-1)}
                    className="rounded px-3 py-2 border hover:bg-gray-50 transition"
                >
                    ← Back
                </button>
            </div>

            {loading && <p>Loading…</p>}
            {error && <p className="text-red-600">Error: {error}</p>}
            {!loading && !error && items.length === 0 && <p>No analyses found for this video.</p>}

            {!loading && !error && items.length > 0 && (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {items.map((a) => (
                        <div key={a.id} className="flex flex-col items-center rounded-lg shadow p-3">
                            <video controls className="w-full mb-2 rounded">
                                <source src={a.url} type="video/mp4" />
                                Your browser does not support the video tag.
                            </video>

                            <div className="w-full text-sm text-gray-600 mb-2">
                                <div>Type: <span className="font-medium">{a.type}</span></div>
                                <div>Video ID: {a.video_id}</div>
                            </div>

                            <div className="flex gap-2 w-full">
                                <a
                                    className="bg-gray-800 text-white px-3 py-2 rounded hover:bg-gray-900 transition w-full text-center"
                                    href={a.url}
                                    download
                                >
                                    Download
                                </a>
                                <a
                                    className="border px-3 py-2 rounded hover:bg-gray-50 transition w-full text-center"
                                    href={a.url}
                                    target="_blank"
                                    rel="noreferrer"
                                >
                                    Open
                                </a>
                            </div>

                            <div className="mt-2 text-xs text-gray-500 w-full break-words">
                                <div>Bucket: {a.bucket}</div>
                                <div>Object: {a.object_key}</div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
