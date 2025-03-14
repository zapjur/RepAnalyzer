import axios, { AxiosRequestConfig } from 'axios';
import { GetTokenSilentlyOptions } from '@auth0/auth0-react';

const apiClient = axios.create({
    baseURL: 'http://localhost:8080',
    headers: {
        'Content-Type': 'application/json',
    },
});

export const setupInterceptors = (
    getAccessTokenSilently: (options?: GetTokenSilentlyOptions) => Promise<string>,
) => {
    apiClient.interceptors.request.use(
        async (config: AxiosRequestConfig) => {
            try {
                const token = await getAccessTokenSilently();
                if (token && config.headers) {
                    config.headers.Authorization = `Bearer ${token}`;
                }
            } catch (error) {
                console.error("Failed to get access token", error);
            }
            return config;
        },
        (error) => Promise.reject(error),
    );
};

export default apiClient;
