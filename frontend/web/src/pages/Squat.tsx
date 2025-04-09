import React, { useEffect } from 'react';
import { useVideos } from '../contexts/VideosContext';
import { useAuth0 } from '@auth0/auth0-react';

const Squat: React.FC = () => {
    const { videos, fetchVideos } = useVideos();
    const { user, isAuthenticated, isLoading } = useAuth0();

    useEffect(() => {
        if (user?.sub) {
            fetchVideos("Squat", user.sub);
        }

    }, []);

    const squatVideos = videos["Squat"] || [];

    return (
        <div>
            <h2 className="text-2xl font-bold mb-4">Squat Videos</h2>
            {squatVideos.length === 0 ? (
                <p>No videos yet</p>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {squatVideos.map((video, idx) => (
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

export default Squat;
