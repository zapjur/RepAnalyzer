import React, { useEffect } from 'react';
import { useAuth0 } from "@auth0/auth0-react";
import { setupInterceptors } from '../api/axios';

const Dashboard: React.FC = () => {
    const { user, logout, getAccessTokenSilently, isAuthenticated } = useAuth0();

    useEffect(() => {
        if (isAuthenticated) {
            setupInterceptors(getAccessTokenSilently);
        }
    }, [isAuthenticated, getAccessTokenSilently]);

    return (
        <div className="flex flex-col items-center gap-4 p-10">
            <h2 className="text-2xl font-bold">Hello, {user?.name}</h2>

            <button
                className="bg-red-500 text-white px-4 py-2 rounded"
                onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
            >
                Wyloguj się
            </button>
        </div>
    );
};

export default Dashboard;
