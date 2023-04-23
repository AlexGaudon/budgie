import { useEffect, useState } from "react";
import { z } from "zod";
import { Transaction, transactionSchema } from "../types";

import { type CreateTransactionForm } from "../components/oldcreatetransaction";

const transactionResponse = z.object({
    data: z.array(transactionSchema),
});

export const useTransactions = () => {
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [refresh, setRefresh] = useState(false);

    const createTransaction = async (data: CreateTransactionForm) => {
        let res = await fetch("/api/transactions", {
            method: "POST",
            body: JSON.stringify({
                type: data.type,
                description: data.description,
                category: data.category,
                amount: Number(
                    Math.round(parseFloat(data.amount) * 100).toString()
                ),
                date: new Date(data.date),
            }),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            let resp = await res.json();
            if ("data" in resp) {
                let parsed = transactionSchema.parse(resp.data);
                setTransactions([...transactions, parsed]);
                console.log(`PARSED: ${JSON.stringify(parsed)}`);
            }
        } else {
            let resp = await res.json();
            console.log(resp);
        }
    };

    const fetchTransactions = async () => {
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
                return parsedData;
            }
        } catch (error) {
            console.log(error);
        }
        return [];
    };

    useEffect(() => {
        const refreshTransactions = async () => {
            let t = await fetchTransactions();
            setTransactions(t);
        };
        refreshTransactions();
    }, [refresh]);

    return {
        transactions,
        refreshTransactions: () => {
            setRefresh((v) => !v);
        },
        createTransaction,
    };
};
