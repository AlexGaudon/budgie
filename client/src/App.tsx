import "./App.css";
import { useAuth } from "./hooks/useAuth";

import { Login } from "./pages/Login";

import LandingPage from "./pages/Landing";

function App() {
    let auth = useAuth();

    if (auth.isLoading) {
        return <p>Loading...</p>;
    }

    if (auth.user != null) {
        return (
            <div>
                <h1>Hello {auth.user.username}</h1>
                <br />
            </div>
        );
    }
    return <LandingPage></LandingPage>;
}

export default App;
