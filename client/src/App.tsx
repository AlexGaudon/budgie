import "./App.css";
import { useAuth } from "./hooks/useAuth";

import LandingPage from "./pages/Landing";
import { AddTransaction } from "./components/AddTransaction";

function App() {
    let auth = useAuth();

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
