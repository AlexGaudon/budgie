import { Navigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

import { TransactionTable } from "../components/TransactionTable";
import { AddTransaction } from "../components/AddTransaction";
import { ImportDropZone } from "../components/ImportDropzone";
import { useState } from "react";

export const Transactions = () => {
    const auth = useAuth();

    const [isAdding, setIsAdding] = useState(false);

    if (!auth.isLoading && !auth.isLoggedIn) {
        return <Navigate to="/"></Navigate>;
    }

    return (
        <div>
            <button
                onClick={() => {
                    console.log("adding");
                    setIsAdding(true);
                    setTimeout(() => {
                        setIsAdding(false);
                    }, 200);
                }}
            >
                Add Transaction
            </button>
            <TransactionTable isAdding={isAdding}></TransactionTable>
        </div>
    );
};
