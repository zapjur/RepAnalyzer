import React, { useEffect } from 'react';
import { useVideos } from '../contexts/VideosContext';
import { useAuth0 } from '@auth0/auth0-react';
import {Link} from "react-router-dom";

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
                        <div key={idx} className="flex flex-col items-center">
                            <video controls className="w-full mb-2">
                                <source src={video.url} type="video/mp4" />
                                Your browser does not support the video tag.
                            </video>
                            <Link
                                className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition inline-block text-center"
                                to={`/analysis/bench/${video.id}`}
                            >
                                View Analysis
                            </Link>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default Bench;
