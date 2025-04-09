import React, { useEffect } from 'react';
import { useVideos } from '../contexts/VideosContext';
import { useAuth0 } from '@auth0/auth0-react';

const Squat: React.FC = () => {
    const { videos, fetchVideos } = useVideos();
    const { user, isAuthenticated, isLoading } = useAuth0();

    useEffect(() => {
        if (user?.sub) {
            fetchVideos("Deadlift", user.sub);
        }

    }, []);

    const deadliftVideos = videos["Deadlift"] || [];

    return (
        <div>
            <h2 className="text-2xl font-bold mb-4">Deadlift Videos</h2>
            {deadliftVideos.length === 0 ? (
                <p>No videos yet</p>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {deadliftVideos.map((video, idx) => (
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
