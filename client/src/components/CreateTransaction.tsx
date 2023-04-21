import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";

/*
{
    "description": "Sobeys2",
    "category": "groceries",
    "amount": 10503,
    "date": "2023-04-15T00:00:00Z"
}
*/

const createTransactionSchema = z.object({
    description: z.string(),
    category: z.string(),
    amount: z.string(),
    date: z.string(),
});

type CreateTransactionForm = z.infer<typeof createTransactionSchema>;

export const CreateTransaction = ({
    onCreateTransaction,
}: {
    onCreateTransaction: () => void;
}) => {
    const {
        register,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm<CreateTransactionForm>({
        resolver: zodResolver(createTransactionSchema),
        defaultValues: {
            date: new Date().toISOString().substring(0, 10),
        },
    });

    const onSubmit: SubmitHandler<CreateTransactionForm> = async (data) => {
        console.log(data);

        let res = await fetch("/api/transactions", {
            method: "POST",
            body: JSON.stringify({
                description: data.description,
                category: data.category,
                amount: ((Number.parseFloat(data.amount) * 50) / 50) * 100,
                date: new Date(data.date),
            }),
            headers: {
                Authorization: window.localStorage.token,
            },
        });

        if (res.ok) {
            let json = await res.json();
            console.log(json);
            onCreateTransaction();
            reset();
        }

        console.log(await res.json());
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
                {errors.amount && errors.amount.message}
                <input
                    autoComplete="off"
                    type="number"
                    min="0"
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

                <input
                    type="submit"
                    className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md cursor-pointer"
                />
            </form>
        </div>
    );
};
