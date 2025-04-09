import React, { createContext, useContext, useState } from 'react';
import apiClient from "../api/axios.ts";

type Video = {
    url: string;
    exercise_name: string;
    uploaded_at: string;
};

type VideosMap = Record<string, Video[]>;

type VideosContextType = {
    videos: VideosMap;
    fetchVideos: (exercise: string, auth0Id: string) => Promise<void>;
};

const VideosContext = createContext<VideosContextType | undefined>(undefined);

export const useVideos = () => {
    const context = useContext(VideosContext);
    if (!context) throw new Error("useVideos must be used within a VideosProvider");
    return context;
};

export const VideosProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [videos, setVideos] = useState<VideosMap>({});

    const fetchVideos = async (exercise: string, auth0Id: string) => {
        console.log(`Fetching videos for ${exercise}...`);
        if (videos[exercise]) return;

        try {
            const encodedExercise = encodeURIComponent(exercise);
            const res = await apiClient.get(`/videos/${auth0Id}/${encodedExercise}`);
            setVideos((prev) => ({ ...prev, [exercise]: res.data }));
        } catch (err) {
            console.error(`Error fetching ${exercise} videos:`, err);
        }
    };

    return (
        <VideosContext.Provider value={{ videos, fetchVideos }}>
            {children}
        </VideosContext.Provider>
    );
};
