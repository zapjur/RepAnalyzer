import React, { useEffect, useState } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import {Link, Outlet} from "react-router-dom";
import apiClient from "../api/axios.ts";

const Dashboard: React.FC = () => {
    const { user, logout, getAccessTokenSilently, isAuthenticated } = useAuth0();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [selectedExercise, setSelectedExercise] = useState("Squat");

    useEffect(() => {
        if (isAuthenticated && user?.sub) {
            const init = async () => {
                if (isAuthenticated) {
                    try {
                        const token = await getAccessTokenSilently();
                        localStorage.setItem('access_token', token);
                    } catch (err) {
                        console.error("Error getting access token", err);
                    }
                }
            };
            init();
            fetchUser();
            console.log("Is Authenticated:", isAuthenticated);
            console.log("User:", user);
        }
    }, [isAuthenticated, getAccessTokenSilently, user]);

    const fetchUser = async () => {
        try {
            const response = await apiClient.get(`/users`);
            console.log('User data:', response.data);
        } catch (error) {
            console.error('Failed to fetch user:', error);
        }
    };

    const handleUpload = async () => {
        if (!selectedFile) return alert("Please select a file first!");

        const formData = new FormData();
        formData.append("file", selectedFile);
        formData.append("exercise", selectedExercise);

        try {
            const response = await apiClient.post(`/upload`, formData, {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
            });
            console.log("Upload success:", response.data);
            setIsModalOpen(false);
        } catch (error) {
            console.error("Upload failed:", error);
            alert("Upload failed!");
        }
    };



    return (
        <div className="flex h-screen">
            <aside className="w-64 bg-stone-50 text-white p-6 h-screen flex flex-col">
                <Link to={"/dashboard"}>
                    <img src="/repanalyzer-logo-small.png" alt="Logo" className="w-32 mx-auto mb-4" />
                </Link>

                <nav className="flex flex-col flex-1">
                    <ul className="space-y-2">
                        <li>
                            <Link to="squat" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Squat</Link>
                        </li>
                        <li>
                            <Link to="bench" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Bench Press</Link>
                        </li>
                        <li>
                            <Link to="deadlift" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Deadlift</Link>
                        </li>
                        <li>
                            <Link to="settings" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Settings</Link>
                        </li>
                    </ul>
                </nav>

                <button
                    onClick={() => setIsModalOpen(true)}
                    className="mb-3 bg-yellow-300 text-stone-950 block p-2 rounded-full hover:bg-yellow-400 mt-auto"
                >
                    Upload Video
                </button>

                <button
                    className="bg-red-500 text-stone-950 block p-2 rounded-full hover:bg-red-600 mt-auto"
                    onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
                >
                    Logout
                </button>
            </aside>


            <main className="flex-1 p-6 bg-stone-200">
                <Outlet/>
            </main>
            {isModalOpen && (
                <div className="fixed inset-0 backdrop-blur-sm bg-black/20 flex items-center justify-center z-50">
                <div className="bg-white rounded-lg shadow-lg p-6 w-full max-w-md">
                        <h2 className="text-xl font-semibold mb-4">Upload a Video</h2>

                        <div className="mb-4">
                            <label className="block text-sm font-medium text-gray-700 mb-1">Exercise</label>
                            <select
                                className="w-full border border-gray-300 rounded px-3 py-2"
                                value={selectedExercise}
                                onChange={(e) => setSelectedExercise(e.target.value)}
                            >
                                <option>Squat</option>
                                <option>Bench Press</option>
                                <option>Deadlift</option>
                            </select>
                        </div>

                        <div className="mb-4">
                            <label className="block text-sm font-medium text-gray-700 mb-1">Choose a file</label>
                            <input
                                type="file"
                                accept="video/*"
                                className="w-full"
                                onChange={(e) => setSelectedFile(e.target.files?.[0] || null)}
                            />
                        </div>

                        <div className="flex justify-end gap-2">
                            <button
                                className="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
                                onClick={() => setIsModalOpen(false)}
                            >
                                Cancel
                            </button>
                            <button
                                onClick={() => handleUpload()}
                                className="px-4 py-2 bg-yellow-400 rounded hover:bg-yellow-500"
                            >
                                Upload
                            </button>
                        </div>
                    </div>
                </div>
            )}

        </div>
    );
};

export default Dashboard;
