import { useEffect, useState } from "react";
import { type Transaction } from "../types";

import { TransactionRow } from "./TransactionRow";
import { AddTransaction } from "./AddTransaction";
import {
    useCreateTransactionMutation,
    useDeleteTransactionMutation,
    useTransactionQuery,
    useUpdateTransactionMutation,
} from "../hooks/useTransactions";

export const TransactionTable = () => {
    const [editingIndex, setEditingIndex] = useState(-1);
    const [creatingTransaction, setCreatingTransaction] = useState(false);

    const { data: transactions, isLoading, error } = useTransactionQuery();
    const updateTransaction = useUpdateTransactionMutation();
    const deleteTransaction = useDeleteTransactionMutation();

    if (isLoading) {
        return <h1>isloading</h1>;
    }

    const onRowEvent = async (
        eventType: "edit" | "delete",
        subject: Transaction
    ) => {
        if (eventType === "edit" || eventType === "delete") {
            if (eventType === "edit") {
                console.log(subject);
                updateTransaction.mutateAsync(subject);
            } else {
                deleteTransaction.mutateAsync(subject.id);
            }
        }
    };

    const showAddTransaction = () => {
        setCreatingTransaction(true);
    };

    return (
        <div className="w-8/12">
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
                        <th className="text-left px-4 py-2">Date</th>
                        <th className="text-left px-4 py-2">Vendor</th>
                        <th className="text-left px-4 py-2">Description</th>
                        <th className="text-left px-4 py-2">Category</th>
                        <th className="text-left px-4 py-2">Amount</th>
                        <th className="text-left px-4 py-2">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {transactions?.map((row, index) => {
                        console.log(row.id);
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
