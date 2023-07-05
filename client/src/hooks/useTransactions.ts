import { useQuery, useMutation, useQueryClient } from "react-query";

import { z } from "zod";

import { type Transaction, transactionSchema } from "../types";
import { CreateTransactionForm } from "../components/AddTransaction";

const fetchTransactions = async (category: string|undefined, period: string|undefined) => {
    const res = await fetch("/api/transactions");

    if (res.ok) {
        const data = await res.json();
        if ("data" in data) {
            let parsed = z.array(transactionSchema).parse(data.data);
            // Apply categorying if a category is provided
            if (category) {
                parsed = parsed.filter((transaction) => transaction.category === category);
            }
            if (period) {
                parsed = parsed.filter((transaction) => transaction.date.toISOString().substring(0, 7) === period)
            }
            return parsed;
        } else {
            throw new Error("Error fetching transactions");
        }
    }
};

export const useCreateTransactionMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (transaction: CreateTransactionForm) => {
        const res = await fetch("/api/transactions", {
            method: "POST",
            body: JSON.stringify({
                type: transaction.type,
                vendor: transaction.vendor,
                description: transaction.description,
                category_id: transaction.category,
                amount: Number(
                    Math.round(
                        parseFloat(transaction.amount.toString()) * 100
                    ).toString()
                ),
                date: new Date(transaction.date),
            }),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            queryClient.invalidateQueries("transactions");
        } else {
            throw new Error("Error creating transaction");
        }
    });
};

export const useUpdateTransactionMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (transaction: Transaction) => {
        const res = await fetch(`/api/transactions/${transaction.id}`, {
            method: "PUT",
            body: JSON.stringify({
                type: transaction.type,
                vendor: transaction.vendor,
                description: transaction.description,
                category_id: transaction.category_id,
                amount: Number(
                    Math.round(
                        parseFloat(transaction.amount.toString()) * 100
                    ).toString()
                ),
                date: new Date(transaction.date),
            }),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            queryClient.invalidateQueries("transactions");
        } else {
            throw new Error("Error updating transaction");
        }
    });
};

export const useDeleteTransactionMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (transactionId: string) => {
        const res = await fetch(`/api/transactions/${transactionId}`, {
            method: "DELETE",
        });
        if (res.ok) {
            queryClient.invalidateQueries("transactions");
        } else {
            throw new Error("Error deleting transaction");
        }
    });
};

export const useTransactionQuery = (category: string | undefined, period: string|undefined) => {
    return useQuery(["transactions", category, period], () =>
        fetchTransactions(category, period)
    );
};
