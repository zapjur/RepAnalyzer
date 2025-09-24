import { Route, Routes } from "react-router-dom";
import Landing from "./pages/Landing";
import Dashboard from "./pages/Dashboard";
import Squat from "./pages/Squat";
import BenchPress from "./pages/Bench";
import Deadlift from "./pages/Deadlift";
import Settings from "./pages/Settings";
import AnalysisDeadliftPage from "./pages/AnalysisDeadliftPage";
import AnalysisSquatPage from "./pages/AnalysisSquatPage";
import AnalysisBenchPage from "./pages/AnalysisBenchPage.tsx";

function App() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Landing />} />
                <Route path="/dashboard" element={<Dashboard />}>
                    <Route path="squat" element={<Squat />} />
                    <Route path="bench" element={<BenchPress />} />
                    <Route path="deadlift" element={<Deadlift />} />
                    <Route path="settings" element={<Settings />} />
                </Route>
                <Route path="analysis/deadlift/:videoId" element={<AnalysisDeadliftPage />} />
                <Route path="analysis/squat/:videoId" element={<AnalysisSquatPage />} />
                <Route path="analysis/bench/:videoId" element={<AnalysisBenchPage />} />
            </Routes>
        </>
    );
}

export default App;
