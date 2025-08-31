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

type AnalysesByVideo = Record<number, VideoAnalysis[]>;

type AnalysesContextType = {
    byVideo: AnalysesByVideo;
    fetchByVideo: (videoId: number) => Promise<VideoAnalysis[]>;
};

const AnalysesContext = createContext<AnalysesContextType | undefined>(undefined);

export const useAnalyses = () => {
    const ctx = useContext(AnalysesContext);
    if (!ctx) throw new Error("useAnalyses must be used within an AnalysesProvider");
    return ctx;
};

export const AnalysesProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [byVideo, setByVideo] = useState<AnalysesByVideo>({});

    const fetchByVideo = async (videoId: number) => {
        if (byVideo[videoId]) return byVideo[videoId];

        const res = await apiClient.get(`/video-analysis/${videoId}`);
        const data = (res.data || []) as VideoAnalysis[];

        setByVideo(prev => ({ ...prev, [videoId]: data }));
        return data;
    };

    return (
        <AnalysesContext.Provider value={{ byVideo, fetchByVideo }}>
            {children}
        </AnalysesContext.Provider>
    );
};
