import React, { useEffect } from 'react';
import { useVideos } from '../contexts/VideosContext';
import { useAuth0 } from '@auth0/auth0-react';

const Bench: React.FC = () => {
    const { videos, fetchVideos } = useVideos();
    const { user, isAuthenticated, isLoading } = useAuth0();

    useEffect(() => {
        if (user?.sub) {
            fetchVideos("Bench_Press", user.sub);
        }

    }, []);

    const benchVideos = videos["Bench_Press"] || [];

    return (
        <div>
            <h2 className="text-2xl font-bold mb-4">Bench Press Videos</h2>
            {benchVideos.length === 0 ? (
                <p>No videos yet</p>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {benchVideos.map((video, idx) => (
                        <video key={idx} controls className="w-full">
                            <source src={video.url} type="video/mp4" />
                            Your browser does not support the video tag.
                        </video>
                    ))}
                </div>
            )}
        </div>
    );
};

export default Bench;
