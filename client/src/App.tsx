import "./App.css";
import { useAuth } from "./hooks/useAuth";

import { Login } from "./pages/Login";

import LandingPage from "./pages/Landing";
import { AddTransaction } from "./components/AddTransaction";

function App() {
    let auth = useAuth();

    if (auth.isLoading) {
        return <p>Loading...</p>;
    }

    if (auth.user != null) {
        return (
            <div>
                <AddTransaction></AddTransaction>
            </div>
        );
    }
    return <LandingPage></LandingPage>;
}

export default App;
