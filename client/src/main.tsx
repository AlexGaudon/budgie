import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";

import {
    createBrowserRouter,
    createRoutesFromElements,
    Route,
    RouterProvider,
} from "react-router-dom";
import { Transactions } from "./pages/Transactions";
import { Layout } from "./Layout";
import { Login } from "./pages/Login";
import { AuthProvider } from "./contexts/AuthContext";

import { Budgets } from "./pages/Budgets";
import { Categories } from "./pages/Categories";
import { QueryClient, QueryClientProvider } from "react-query";
import { NewTransaction } from "./pages/NewTransaction";

const router = createBrowserRouter(
    createRoutesFromElements(
        <Route path="/" element={<Layout />}>
            <Route index element={<App />} />
            <Route path="login" element={<Login />} />
            <Route path="categories" element={<Categories />} />
            <Route path="budgets" element={<Budgets />} />
            <Route path="transactions" element={<Transactions />} />
            <Route path="newtransaction" element={<NewTransaction />} />
        </Route>
    )
);

const queryClient = new QueryClient();

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <React.StrictMode>
        <QueryClientProvider client={queryClient}>
            <AuthProvider>
                <RouterProvider router={router} />
            </AuthProvider>
        </QueryClientProvider>
    </React.StrictMode>
);
