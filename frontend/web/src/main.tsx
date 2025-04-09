import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App.tsx";
import AuthProvider from "./providers/auth.tsx";
import "./index.css";
import {VideosProvider} from "./contexts/VideosContext.tsx";

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <BrowserRouter>
            <AuthProvider>
                <VideosProvider>
                    <App />
                </VideosProvider>
            </AuthProvider>
        </BrowserRouter>
    </React.StrictMode>
);
