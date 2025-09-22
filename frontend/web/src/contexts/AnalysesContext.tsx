import React, { createContext, useContext, useState } from "react";
import apiClient from "../api/axios";

export type VideoAnalysis = {
    id: number;
    bucket: string;
    object_key: string;
    type: string;
    url: string;
    csv_url: string | null;
    video_id: number;
};

export type AnalysisPayload = unknown;

type AnalysesByVideo = Record<number, VideoAnalysis[]>;
type AnalysisJSONByVideo = Record<number, AnalysisPayload>;

type FetchResult = {
    videos: VideoAnalysis[];
    analysis: AnalysisPayload;
};

type AnalysesContextType = {
    byVideo: AnalysesByVideo;
    analysisByVideo: AnalysisJSONByVideo;
    fetchByVideo: (videoId: number) => Promise<FetchResult>;
};

const AnalysesContext = createContext<AnalysesContextType | undefined>(undefined);

export const useAnalyses = () => {
    const ctx = useContext(AnalysesContext);
    if (!ctx) throw new Error("useAnalyses must be used within an AnalysesProvider");
    return ctx;
};

export const AnalysesProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [byVideo, setByVideo] = useState<AnalysesByVideo>({});
    const [analysisByVideo, setAnalysisByVideo] = useState<AnalysisJSONByVideo>({});

    const fetchByVideo = async (videoId: number): Promise<FetchResult> => {
        if (byVideo[videoId] && analysisByVideo[videoId]) {
            return { videos: byVideo[videoId], analysis: analysisByVideo[videoId] };
        }

        const res = await apiClient.get(`/video-analysis/${videoId}`);
        const data = res.data as { videos?: VideoAnalysis[]; analysis?: AnalysisPayload };

        const videos = data.videos ?? [];
        const analysis = data.analysis ?? null;

        setByVideo(prev => ({ ...prev, [videoId]: videos }));
        setAnalysisByVideo(prev => ({ ...prev, [videoId]: analysis }));

        return { videos, analysis };
    };

    return (
        <AnalysesContext.Provider value={{ byVideo, analysisByVideo, fetchByVideo }}>
            {children}
        </AnalysesContext.Provider>
    );
};
