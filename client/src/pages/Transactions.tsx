import { Navigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

import { TransactionTable } from "../components/TransactionTable";

export const Transactions = () => {
    const auth = useAuth();

    if (!auth.isLoading && !auth.isLoggedIn) {
        return <Navigate to="/"></Navigate>;
    }

    return <TransactionTable></TransactionTable>;
};
