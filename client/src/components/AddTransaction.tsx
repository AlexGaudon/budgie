import { useEffect } from "react";

import { z } from "zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useCategoriesQuery } from "../hooks/useCategories";
import { useCreateTransactionMutation } from "../hooks/useTransactions";

const createTransactionSchema = z.object({
    vendor: z.string(),
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

export const AddTransaction = ({ onFinish }: { onFinish: () => void }) => {
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

    const {
        data: categories,
        isLoading: categoryLoading,
        error: categoryError,
    } = useCategoriesQuery();

    const createTransactionMutation = useCreateTransactionMutation();

    if (categoryLoading) {
        return <h1>Loading...</h1>;
    }

    const onSubmit: SubmitHandler<CreateTransactionForm> = async (data) => {
        let input = {
            vendor: data.vendor,
            description: data.description,
            category: data.category,
            amount: data.amount,
            type: data.type,
            date: data.date,
        };

        await createTransactionMutation.mutateAsync(input);

        onFinish();
    };

    return (
        <div className="flex items-center justify-center h-full w-8/12">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div className="grid gap-4">
                    <label htmlFor="vendor">Vendor</label>
                    <input
                        id="vendor"
                        autoComplete="off"
                        type="text"
                        placeholder="Walmart"
                        {...register("vendor", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />

                    <label htmlFor="description">Description</label>
                    <input
                        id="description"
                        autoComplete="off"
                        type="text"
                        placeholder="Snacks"
                        {...register("description", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />

                    <label htmlFor="category">Category</label>
                    <select
                        id="category"
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        {...register("category", {
                            required: true,
                            maxLength: 80,
                        })}
                    >
                        {categories?.map((c) => {
                            return (
                                <option key={c.id} value={c.id}>
                                    {c.name}
                                </option>
                            );
                        })}
                    </select>

                    <label htmlFor="amount">Amount</label>
                    <input
                        id="amount"
                        autoComplete="off"
                        type="number"
                        step="0.01"
                        placeholder="0.00"
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
