import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";

import { useTransactions } from "../hooks/useTransactions";

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

export const CreateTransaction = () => {
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

    const { createTransaction, refreshTransactions } = useTransactions();

    const onSubmit: SubmitHandler<CreateTransactionForm> = async (data) => {
        await createTransaction(data);

        refreshTransactions();
    };

    return (
        <div className="block">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <input
                    autoComplete="off"
                    type="text"
                    placeholder="Description"
                    {...register("description", {
                        required: true,
                        maxLength: 80,
                    })}
                    className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />

                <input
                    autoComplete="off"
                    type="text"
                    placeholder="Category"
                    {...register("category", {
                        required: true,
                        maxLength: 80,
                    })}
                    className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />

                <input
                    autoComplete="off"
                    type="number"
                    step="0.01"
                    placeholder="Amount"
                    {...register("amount", {
                        required: true,
                        maxLength: 80,
                    })}
                    className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />

                <input
                    autoComplete="off"
                    type="date"
                    placeholder="Date"
                    {...register("date", {
                        required: true,
                        maxLength: 80,
                    })}
                    className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />

                <label htmlFor="field-expense">
                    <input
                        {...register("type")}
                        type="radio"
                        value="expense"
                        id="field-expense"
                    />
                    Expense
                </label>
                <label htmlFor="field-income">
                    <input
                        {...register("type")}
                        type="radio"
                        value="income"
                        id="field-income"
                    />
                    Income
                </label>

                <input
                    type="submit"
                    className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md cursor:pointer"
                />
            </form>
        </div>
    );
};
