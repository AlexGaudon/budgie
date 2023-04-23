import { create } from "zustand";
import { transactionSchema, Transaction } from "./types";
import { z } from "zod";

import { CreateTransactionForm } from "./components/AddTransaction";

interface TransactionStore {
    transactions: Transaction[];
    isLoading: boolean;
    error: string | null;
    fetchTransactions: () => Promise<void>;
    deleteTransaction: (transaction: Transaction) => Promise<void>;
    createTransaction: (transaction: CreateTransactionForm) => Promise<void>;
    updateTransaction: (transaction: Transaction) => Promise<void>;
}

export const useTransactionStore = create<TransactionStore>((set) => ({
    transactions: [],
    isLoading: false,
    error: null,
    createTransaction: async (transaction: CreateTransactionForm) => {
        try {
            const res = await fetch("/api/transactions", {
                method: "POST",
                body: JSON.stringify({
                    type: transaction.type,
                    description: transaction.description,
                    category: transaction.category,
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
            const data = await res.json();

            if (res.ok && "data" in data) {
                let parsed = transactionSchema.parse(data.data);

                set((state) => ({
                    transactions: [parsed, ...state.transactions],
                }));
            }
        } catch (error) {
            console.log(error);
            set({
                error: "Failed to create transaction",
            });
        }
    },
    deleteTransaction: async (transaction: Transaction) => {
        try {
            await fetch(`/api/transactions/${transaction.id}`, {
                method: "DELETE",
            });

            set((state) => {
                const index = state.transactions.findIndex(
                    (x) => x.id === transaction.id
                );
                if (index !== -1) {
                    const newTransactions = [
                        ...state.transactions.slice(0, index),
                        ...state.transactions.slice(index + 1),
                    ];
                    return {
                        transactions: newTransactions,
                    };
                } else {
                    return state;
                }
            });
        } catch (error) {
            console.log(error);
            set({
                error: "Failed to delete transaction",
            });
        }
    },

    updateTransaction: async (transaction: Transaction) => {
        try {
            const res = await fetch(`/api/transactions/${transaction.id}`, {
                method: "PUT",
                body: JSON.stringify({
                    type: transaction.type,
                    description: transaction.description,
                    category: transaction.category,
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
                let body = await res.json();
                if ("data" in body) {
                    let parsed = transactionSchema.parse(body.data);
                    set((state) => {
                        let theTransaction = state.transactions.find(
                            (x) => x.id === transaction.id
                        );

                        let idx = state.transactions.findIndex(
                            (x) => x.id === transaction.id
                        );
                        if (~idx) {
                            theTransaction = { ...parsed };

                            state.transactions[idx] = theTransaction;

                            return {
                                transactions: state.transactions,
                            };
                        }
                        return {};
                    });
                }
            }
        } catch (error) {
            set({
                error: "Failed to update transaction",
            });
        }
    },
    fetchTransactions: async () => {
        set({ isLoading: true });

        try {
            const response = await fetch("/api/transactions");
            const data = await response.json();
            if ("data" in data) {
                // Validate fetched transactions using Zod schema
                const parsedTransactions = z
                    .array(transactionSchema)
                    .parse(data.data);
                set({ transactions: parsedTransactions, error: null });
            }
        } catch (error) {
            console.log(error);
            set({
                error: "Failed to fetch transactions",
                transactions: [],
                isLoading: false,
            });
        } finally {
            set({ isLoading: false });
        }
    },
}));
