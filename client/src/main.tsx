import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";

import {
    createBrowserRouter,
    createRoutesFromElements,
    Route,
    RouterProvider,
    useNavigate,
} from "react-router-dom";
import { Transactions } from "./pages/Transactions";
import { Layout } from "./Layout";
import { Login } from "./pages/Login";
import { AuthProvider } from "./contexts/AuthContext";

import { Budgets } from "./pages/Budgets";

const router = createBrowserRouter(
    createRoutesFromElements(
        <Route path="/" element={<Layout />}>
            <Route index element={<App />} />
            <Route path="login" element={<Login />} />
            <Route path="budgets" element={<Budgets />} />
            <Route path="transactions" element={<Transactions />} />
        </Route>
    )
);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <React.StrictMode>
        <AuthProvider>
            <RouterProvider router={router} />
        </AuthProvider>
    </React.StrictMode>
);
