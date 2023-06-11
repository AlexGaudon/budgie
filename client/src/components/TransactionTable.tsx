import { useEffect, useState } from "react";
import { Category, type Transaction } from "../types";

import { useLocation } from "react-router-dom";

import { TransactionRow } from "./TransactionRow";
import {
    useCreateTransactionMutation,
    useDeleteTransactionMutation,
    useTransactionQuery,
    useUpdateTransactionMutation,
} from "../hooks/useTransactions";
import { useCategoriesQuery } from "../hooks/useCategories";

type TransactionTableProps = {
    isAdding: boolean;
};

export const TransactionTable = ({ isAdding }: TransactionTableProps) => {
    const [editingIndex, setEditingIndex] = useState(0);

    const { data: categories } = useCategoriesQuery();

    const createTransactionMutation = useCreateTransactionMutation();

    const updateTransaction = useUpdateTransactionMutation();
    const deleteTransaction = useDeleteTransactionMutation();

    const location = useLocation();
    const searchParams = new URLSearchParams(location.search);
    const filter = searchParams.get("filter") as string | undefined;

    const { data: transactions, isLoading } = useTransactionQuery(filter);

    useEffect(() => {
        if (isAdding) {
            let defaultCategory = categories?.find(
                (e: Category) => e.name === "Uncategorized"
            );

            async function makeDefault() {
                if (defaultCategory === undefined) {
                    return;
                }
                await createTransactionMutation.mutateAsync({
                    vendor: "",
                    description: "",
                    category: defaultCategory.id,
                    amount: "0.00",
                    type: "expense",
                    date: new Date().toISOString().substring(0, 10),
                });
                transactions?.unshift();
            }

            makeDefault();
        }
    }, [isAdding]);

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

    return (
        <div className="w-8/12">
            <div className="container"></div>
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
