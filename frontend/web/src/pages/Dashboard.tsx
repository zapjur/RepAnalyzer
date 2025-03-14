import React, { useEffect } from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import apiClient, { setupInterceptors } from '../api/axios';

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
        <div className="flex flex-col items-center gap-4 p-10">
            <h1 className="text-2xl">Hello, {user?.name}</h1>
            <button
                className="bg-red-500 text-white px-4 py-2 rounded"
                onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
            >
                Logout
            </button>
        </div>
    );
};

export default Dashboard;
