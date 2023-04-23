import { useState } from "react";
import { useTransactionStore } from "../store";

import { type NewTransaction } from "../types";

import { z } from "zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

type NewTransactionProps = {
    onFinish: () => void;
};

const createTransactionSchema = z.object({
    description: z.string(),
    category: z.string(),
    amount: z.string(),
    date: z.string(),
    type: z.string().refine(
        (str) => {
            return str == "income" || str == "expense";
        },
        {
            message: 'Type must be equal to "income" or "expense"',
        }
    ),
});

export type CreateTransactionForm = z.infer<typeof createTransactionSchema>;

export const AddTransaction = ({ onFinish }: NewTransactionProps) => {
    const {
        register,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm<CreateTransactionForm>({
        resolver: zodResolver(createTransactionSchema),
        defaultValues: {
            date: new Date().toISOString().substring(0, 10),
            type: "expense",
        },
    });

    const transactionStore = useTransactionStore();

    const onSubmit: SubmitHandler<CreateTransactionForm> = async (data) => {
        transactionStore.createTransaction({
            description: data.description,
            category: data.category,
            amount: data.amount,
            type: data.type,
            date: data.date,
        });
    };

    return (
        <div className="flex items-center justify-center h-full">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                    <label htmlFor="description">Description</label>
                    <input
                        id="description"
                        autoComplete="off"
                        type="text"
                        placeholder="Walmart"
                        {...register("description", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />

                    <label htmlFor="category">Category</label>
                    <input
                        id="category"
                        autoComplete="off"
                        type="text"
                        placeholder="Groceries"
                        {...register("category", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />

                    <label htmlFor="amount">Amount</label>
                    <input
                        id="amount"
                        autoComplete="off"
                        type="number"
                        step="0.01"
                        placeholder="159.40"
                        {...register("amount", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />

                    <label htmlFor="date">Date of Transaction</label>
                    <input
                        id="date"
                        autoComplete="off"
                        type="date"
                        placeholder="Date"
                        {...register("date", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    <label>Transaction Type</label>
                    <div>
                        <label htmlFor="field-expense" className="mr-2">
                            <input
                                className="mr-1"
                                {...register("type")}
                                type="radio"
                                value="expense"
                                id="field-expense"
                            />
                            Expense
                        </label>

                        <label htmlFor="field-income">
                            <input
                                className="mr-1"
                                {...register("type")}
                                type="radio"
                                value="income"
                                id="field-income"
                            />
                            Income
                        </label>
                    </div>
                </div>

                <input
                    type="submit"
                    className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md cursor:pointer"
                />
            </form>
        </div>
    );
};
