import { useAuth0 } from "@auth0/auth0-react";

export default function Dashboard() {
    const { user, logout } = useAuth0();

    return (
        <div className="flex flex-col items-center gap-4 p-10">
            <h1 className="text-2xl">Hello</h1>
            <button
                className="bg-red-500 text-white px-4 py-2 rounded"
                onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
            >
                Logout
            </button>
        </div>
    );
}
