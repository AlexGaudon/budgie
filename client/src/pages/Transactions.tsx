import React, { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { CreateTransaction } from "../components/oldcreatetransaction";
import { useAuth } from "../hooks/useAuth";

import { useTransactionStore } from "../store";

export const Transactions = () => {
    const auth = useAuth();

    const transactionStore = useTransactionStore();

    const fetchTransactions = useTransactionStore(
        (state) => state.fetchTransactions
    );

    useEffect(() => {
        fetchTransactions();
    }, [fetchTransactions]);

    if (transactionStore.error) {
        return <h1>{transactionStore.error}</h1>;
    }

    if (transactionStore.isLoading) {
        return <h1>{transactionStore.isLoading}</h1>;
    }

    if (!auth.isLoading && !auth.isLoggedIn) {
        return <Navigate to="/"></Navigate>;
    }

    // const deleteTransaction = async (id: string) => {
    //     let res = await fetch(`/api/transactions/${id}`, {
    //         method: "DELETE",
    //     });
    //     if (res.ok) {
    //         refreshTransactions();
    //     }
    // };

    return (
        <div className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold mb-6">Transactions</h1>
            <CreateTransaction />
            <div className="mt-4">
                {transactionStore.transactions.map((t) => (
                    <div
                        key={t.id}
                        className="border border-gray-300 rounded-md p-4 mb-4 flex justify-between items-center"
                    >
                        <div>
                            <p className="text-gray-600 text-sm">ID: {t.id}</p>
                            <p>
                                {t.description} for {t.amount} in {t.category}
                            </p>
                        </div>
                        <button
                            onClick={() => {
                                transactionStore.deleteTransaction(t);
                            }}
                            className="px-4 py-2 bg-red-500 text-white rounded-md"
                        >
                            Delete
                        </button>
                    </div>
                ))}
            </div>
        </div>
    );
};
