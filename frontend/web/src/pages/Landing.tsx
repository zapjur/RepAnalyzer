import LoginButton from "../components/LoginButton";

export default function Landing() {
    return (
        <div className="flex flex-col items-center justify-center h-screen text-center">
            <h1 className="text-4xl font-bold mb-6">RepAnalyzer</h1>
            <LoginButton />
        </div>
    );
}
