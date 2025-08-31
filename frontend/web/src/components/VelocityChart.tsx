import React, { useEffect, useMemo, useRef, useState } from "react";
import {
    ComposedChart,
    Line,
    XAxis,
    YAxis,
    Tooltip,
    CartesianGrid,
    ReferenceLine,
    ReferenceDot,
    ResponsiveContainer,
} from "recharts";

export interface VelocityChartProps {
    csvUrl: string;
    videoRef: React.RefObject<HTMLVideoElement | null>;
    useSmoothedVelocity?: boolean;
    height?: number;
    timeOffsetSeconds?: number;
    seekBehavior?: "video" | "chart"; // currently unused
}

interface Row {
    frame: number;
    t: number;
    x_px: number;
    y_px: number;
    vy_m_s: number;
    vy_smooth_m_s: number;
    meters_per_pixel: number;
}

async function fetchCsv(url: string): Promise<string> {
    const res = await fetch(url, { credentials: "omit" });
    if (!res.ok) throw new Error(`CSV download failed: ${res.status} ${res.statusText}`);
    return await res.text();
}

function parseCsv(text: string): Row[] {
    const lines = text.trim().split(/\r?\n/);
    if (lines.length <= 1) return [];
    const header = lines[0].split(",").map((s) => s.trim());
    const idxOf = (name: string) => header.findIndex((h) => h === name);

    const idx = {
        frame: idxOf("frame"),
        t: idxOf("t"),
        x_px: idxOf("x_px"),
        y_px: idxOf("y_px"),
        vy_m_s: idxOf("vy_m_s"),
        vy_smooth_m_s: idxOf("vy_smooth_m_s"),
        meters_per_pixel: idxOf("meters_per_pixel"),
    } as const;

    const rows: Row[] = [];
    for (let i = 1; i < lines.length; i++) {
        const parts = lines[i].split(",").map((s) => s.trim());
        const n = (j: number) => (j >= 0 && j < parts.length ? Number(parts[j]) : NaN);
        rows.push({
            frame: n(idx.frame),
            t: n(idx.t),
            x_px: n(idx.x_px),
            y_px: n(idx.y_px),
            vy_m_s: n(idx.vy_m_s),
            vy_smooth_m_s: n(idx.vy_smooth_m_s),
            meters_per_pixel: n(idx.meters_per_pixel),
        });
    }
    return rows.filter((r) => Number.isFinite(r.t)).sort((a, b) => a.t - b.t);
}

function interpolateY(data: { t: number; y: number }[], t: number): number | null {
    if (!data.length) return null;
    if (t <= data[0].t) return data[0].y;
    if (t >= data[data.length - 1].t) return data[data.length - 1].y;
    let lo = 0,
        hi = data.length - 1;
    while (lo <= hi) {
        const mid = (lo + hi) >> 1;
        if (data[mid].t < t) lo = mid + 1;
        else hi = mid - 1;
    }
    const i = Math.max(1, lo);
    const a = data[i - 1];
    const b = data[i];
    const alpha = (t - a.t) / (b.t - a.t);
    return a.y + alpha * (b.y - a.y);
}

const VelocityChart: React.FC<VelocityChartProps> = ({
                                                         csvUrl,
                                                         videoRef,
                                                         useSmoothedVelocity = true,
                                                         height = 200,
                                                         timeOffsetSeconds = 0,
                                                     }) => {
    const [csvText, setCsvText] = useState<string>("");
    const [error, setError] = useState<string | null>(null);
    const [videoTime, setVideoTime] = useState<number>(0);

    const rafRef = useRef<number | undefined>(undefined);

    useEffect(() => {
        let alive = true;
        setCsvText("");
        setError(null);
        fetchCsv(csvUrl)
            .then((t) => alive && setCsvText(t))
            .catch((e) => alive && setError(e.message || String(e)));
        return () => {
            alive = false;
        };
    }, [csvUrl]);

    const rows = useMemo(() => (csvText ? parseCsv(csvText) : []), [csvText]);

    const series = useMemo(() => {
        const key = useSmoothedVelocity ? "vy_smooth_m_s" : "vy_m_s";
        return rows
            .map((r) => ({ t: r.t, y: (r as any)[key] as number }))
            .filter((p) => Number.isFinite(p.t) && Number.isFinite(p.y));
    }, [rows, useSmoothedVelocity]);

    const [tMin, tMax, yMin, yMax] = useMemo((): [number, number, number, number] => {
        if (!series.length) return [0, 0, -1, 1];
        const t0 = series[0].t;
        const t1 = series[series.length - 1].t;
        let min = series[0].y,
            max = series[0].y;
        for (let i = 1; i < series.length; i++) {
            const y = series[i].y;
            if (y < min) min = y;
            if (y > max) max = y;
        }
        const pad = Math.max(0.05 * (max - min), 0.05);
        return [t0, t1, min - pad, max + pad];
    }, [series]);

    const startLoop = () => {
        const v = videoRef.current;
        if (!v) return;
        const step = () => {
            setVideoTime(v.currentTime || 0);
            if (!v.paused && !v.ended) {
                rafRef.current = requestAnimationFrame(step);
            }
        };
        rafRef.current = requestAnimationFrame(step);
    };

    const stopLoop = () => {
        const id = rafRef.current;
        if (typeof id === "number") cancelAnimationFrame(id);
        rafRef.current = undefined;
    };

    useEffect(() => {
        const v = videoRef.current;
        if (!v) return;
        const onPlay = () => {
            stopLoop();
            startLoop();
        };
        const onPauseOrEnded = () => {
            stopLoop();
            setVideoTime(v.currentTime || 0);
        };
        const onTimeUpdate = () => setVideoTime(v.currentTime || 0);
        const onSeeked = () => setVideoTime(v.currentTime || 0);

        v.addEventListener("play", onPlay);
        v.addEventListener("pause", onPauseOrEnded);
        v.addEventListener("ended", onPauseOrEnded);
        v.addEventListener("timeupdate", onTimeUpdate);
        v.addEventListener("seeked", onSeeked);

        return () => {
            v.removeEventListener("play", onPlay);
            v.removeEventListener("pause", onPauseOrEnded);
            v.removeEventListener("ended", onPauseOrEnded);
            v.removeEventListener("timeupdate", onTimeUpdate);
            v.removeEventListener("seeked", onSeeked);
            stopLoop();
        };
    }, [videoRef]);

    const chartTime = useMemo(() => videoTime + timeOffsetSeconds, [videoTime, timeOffsetSeconds]);
    const clampedChartTime = Math.max(tMin, Math.min(tMax, chartTime));
    const dotY = useMemo(() => interpolateY(series, clampedChartTime), [series, clampedChartTime]);

    return (
        <div className="w-full">
            <div className="flex items-center justify-between mb-2">
                <h3 className="text-base font-semibold text-stone-900">Vertical velocity (m/s)</h3>
                <div className="text-xs opacity-70 tabular-nums">
                    t(video) = {videoTime.toFixed(2)} s, t(chart) = {clampedChartTime.toFixed(2)} s
                    {dotY != null ? `, v = ${dotY.toFixed(3)} m/s` : ""}
                </div>
            </div>

            <div style={{ height, pointerEvents: "none" }}>
                <ResponsiveContainer width="100%" height="100%">
                    <ComposedChart data={series} margin={{ top: 10, right: 20, bottom: 20, left: 0 }}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis
                            dataKey="t"
                            type="number"
                            domain={[tMin, tMax]}
                            allowDataOverflow
                            tickFormatter={(v: number) => `${v.toFixed(1)}s`}
                            interval="preserveEnd"
                        />
                        <YAxis dataKey="y" type="number" domain={[yMin, yMax]} tickFormatter={(v: number) => v.toFixed(2)} />
                        <Tooltip
                            isAnimationActive={false}
                            formatter={(value: unknown, name: string) => [
                                typeof value === "number" ? value.toFixed(3) : String(value),
                                name === "y" ? "v (m/s)" : name,
                            ]}
                            labelFormatter={(label: unknown) => `t=${Number(label).toFixed(3)} s`}
                        />
                        <Line type="monotone" dataKey="y" dot={false} strokeWidth={2} isAnimationActive={false} />
                        {Number.isFinite(clampedChartTime) && (
                            <ReferenceLine
                                x={clampedChartTime}
                                strokeDasharray="4 4"
                                isFront
                                ifOverflow="hidden"
                                isAnimationActive={false}
                            />
                        )}
                        {dotY != null && Number.isFinite(clampedChartTime) && (
                            <ReferenceDot x={clampedChartTime} y={dotY} r={4} isFront ifOverflow="hidden" isAnimationActive={false} />
                        )}
                    </ComposedChart>
                </ResponsiveContainer>
            </div>

            {error && <div className="text-red-600 text-xs mt-2">CSV error: {error}</div>}
            <p className="text-[11px] mt-2 opacity-60">Chart is synced with the video playback.</p>
        </div>
    );
};

export default VelocityChart;
