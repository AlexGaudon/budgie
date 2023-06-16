import { useQuery, useMutation, useQueryClient } from "react-query";

import { z } from "zod";

import { type Transaction, transactionSchema } from "../types";
import { CreateTransactionForm } from "../components/AddTransaction";
import { useCategoriesQuery } from "./useCategories";

const fetchTransactions = async (filter: string|undefined) => {
    const res = await fetch("/api/transactions");

    if (res.ok) {
        const data = await res.json();
        if ("data" in data) {
            const parsed = z.array(transactionSchema).parse(data.data);
            // Apply filtering if a filter is provided
            if (filter) {
                return parsed.filter((transaction) => transaction.category === filter);
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

    let {data: categories}= useCategoriesQuery();
    let income = categories?.find(x => x.name === 'Income');

    return useMutation(async (transaction: Transaction) => {
        console.log('category: ' + transaction.category);
        console.log('category: ' + transaction.category_id);
        if (transaction.category_id == income?.id) {
            transaction.type='income';
        }
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

export const useTransactionQuery = (filter: string | undefined) => {
    return useQuery(["transactions", filter], () =>
        fetchTransactions(filter)
    );
};
