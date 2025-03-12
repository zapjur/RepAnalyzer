import { useAuth0 } from "@auth0/auth0-react";

export default function LoginButton() {
    const { loginWithRedirect } = useAuth0();

    return (
        <button
            className="bg-blue-500 text-white px-4 py-2 rounded"
            onClick={() =>
                loginWithRedirect({
                    appState: { returnTo: "/dashboard" }
                })
            }
        >
            Login
        </button>
    );
}
