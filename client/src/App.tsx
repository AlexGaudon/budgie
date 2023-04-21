import "./App.css";
import { useAuth } from "./hooks/useAuth";

import { Login } from "./pages/Login";

function App() {
    let auth = useAuth();

    if (auth.user != null) {
        return (
            <div>
                <h1>Hello {auth.user.username}</h1>

                <h1>DASHBOARD GOES HERE</h1>
            </div>
        );
    }

    return <h1>HOME PAGE</h1>;
}

export default App;
