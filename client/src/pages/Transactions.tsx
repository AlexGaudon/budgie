import { useEffect, useState } from "react";
import { type Transaction, transactionSchema } from "../types";

import { z } from "zod";
import { Navigate } from "react-router-dom";

import { CreateTransaction } from "../components/CreateTransaction";
import { useAuth } from "../hooks/useAuth";

const transactionResponse = z.object({
    data: z.array(transactionSchema),
});

export const Transactions = () => {
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [refresh, setRefresh] = useState(false);

    let auth = useAuth();

    if (!auth.isLoggedIn && !auth.isLoading) {
        return <Navigate to="/login" />;
    }

    useEffect(() => {
        async function getTransactions() {
            try {
                const res = await fetch("/api/transactions", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });
                if (res.ok) {
                    const data = await res.json();
                    const parsedData = transactionResponse.parse(data).data;
                    setTransactions(parsedData);
                }
            } catch (error) {
                console.log(error);
            }
        }

        getTransactions();
    }, [refresh]);

    const deleteTransaction = async (id: string) => {
        let res = await fetch(`http://localhost:3000/api/transactions/${id}`, {
            method: "DELETE",
        });
        if (res.ok) {
            setRefresh((v) => !v);
        }
    };

    return (
        <div>
            <CreateTransaction
                onCreateTransaction={() => {
                    setRefresh((v) => !v);
                }}
            ></CreateTransaction>

            {transactions.map((t) => {
                return (
                    <div key={t.id}>
                        <p>({t.id})</p>
                        <p>
                            {t.description} for {t.amount} in {t.category}
                        </p>
                        <button
                            onClick={() => {
                                deleteTransaction(t.id);
                            }}
                        >
                            delete
                        </button>
                        <br />
                    </div>
                );
            })}
        </div>
    );
};
