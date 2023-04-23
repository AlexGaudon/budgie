import { useEffect, useState } from "react";
import { type Transaction } from "../types";

import { useTransactionStore } from "../store";
import { TransactionRow } from "./TransactionRow";
import { AddTransaction } from "./AddTransaction";

export const Transactions = () => {
    const [editingIndex, setEditingIndex] = useState(-1);
    const transactionStore = useTransactionStore();

    const [creatingTransaction, setCreatingTransaction] = useState(false);

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

    const onRowEvent = async (
        eventType: "edit" | "delete",
        subject: Transaction
    ) => {
        if (eventType === "edit" || eventType === "delete") {
            if (eventType === "edit") {
                transactionStore.updateTransaction(subject);
            } else {
                transactionStore.deleteTransaction(subject);
            }
        }
    };

    const showAddTransaction = () => {
        setCreatingTransaction(true);
    };

    return (
        <div>
            <div className="container m-8">
                {(creatingTransaction && (
                    <AddTransaction
                        onFinish={() => {
                            setCreatingTransaction(false);
                        }}
                    />
                )) || (
                    <button
                        className="bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium"
                        onClick={showAddTransaction}
                    >
                        Add Transaction
                    </button>
                )}
            </div>
            <table className="table-auto w-full">
                <thead>
                    <tr className="border">
                        <th className="text-left px-4 py-2">Id</th>
                        <th className="text-left px-4 py-2">Date</th>
                        <th className="text-left px-4 py-2">Description</th>
                        <th className="text-left px-4 py-2">Category</th>
                        <th className="text-left px-4 py-2">Amount</th>
                        <th className="text-left px-4 py-2">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {transactionStore.transactions.map((row, index) => {
                        return (
                            <TransactionRow
                                key={row.id}
                                index={index}
                                row={row}
                                onRowEvent={onRowEvent}
                                editing={index === editingIndex}
                                setEditingIndex={setEditingIndex}
                            ></TransactionRow>
                        );
                    })}
                </tbody>
            </table>
        </div>
    );
};
