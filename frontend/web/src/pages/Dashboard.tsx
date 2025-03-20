import React, { useEffect } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import apiClient, { setupInterceptors } from '../api/axios';
import { Link } from "react-router-dom";

const Dashboard: React.FC = () => {
    const { user, logout, getAccessTokenSilently, isAuthenticated } = useAuth0();

    useEffect(() => {
        if (isAuthenticated && user?.sub) {
            setupInterceptors(getAccessTokenSilently);
            fetchUser(user.sub);
            console.log("Is Authenticated:", isAuthenticated);
            console.log("User:", user);
        }
    }, [isAuthenticated, getAccessTokenSilently, user]);

    const fetchUser = async (auth0Id: string) => {
        try {
            const response = await apiClient.get(`/users/${auth0Id}`);
            console.log('User data:', response.data);
        } catch (error) {
            console.error('Failed to fetch user:', error);
        }
    };


    return (
        <div className="flex h-screen">
            <aside className="w-64 bg-stone-50 text-white p-6 h-screen flex flex-col">
                <Link to={"/dashboard"}>
                    <img src="public/repanalyzer-logo-small.png" alt="Logo" className="w-32 mx-auto mb-4" />
                </Link>

                <nav className="flex flex-col flex-1">
                    <ul className="space-y-2">
                        <li>
                            <a href="#" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Squat</a>
                        </li><li>
                            <a href="#" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Bench Press</a>
                        </li><li>
                            <a href="#" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Deadlift</a>
                        </li>
                        <li>
                            <a href="#" className="text-stone-950 block p-2 rounded hover:bg-stone-300">Settings</a>
                        </li>
                    </ul>
                </nav>

                <button
                    className="bg-yellow-300 text-stone-950 block p-2 rounded-full hover:bg-yellow-400 mt-auto"
                    onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
                >
                    Logout
                </button>
            </aside>


            <main className="flex-1 p-6 bg-gray-100">
                <h1 className="text-2xl font-semibold">Hello, {user?.name}</h1>
                <div className="mt-4 p-6 bg-white rounded shadow-md min-h-[400px]">

                </div>
            </main>
        </div>
    );
};

export default Dashboard;
