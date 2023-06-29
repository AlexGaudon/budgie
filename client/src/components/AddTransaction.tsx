import { useEffect, useState } from "react";

import { z } from "zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useCategoriesQuery } from "../hooks/useCategories";
import {
    useCreateTransactionMutation,
    useTransactionQuery,
} from "../hooks/useTransactions";

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

export const AddTransaction = () => {
    const {
        register,
        handleSubmit,
        formState: { errors },
        setValue,
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
        reset();
    };

    const { data: transactions } = useTransactionQuery(undefined);

    const [vendorOptions, setVendorOptions] = useState<string[]>([]);

    useEffect(() => {
        let vendorOpts: string[] = [];

        if (transactions === undefined) {
            setVendorOptions(vendorOpts);
            return;
        }

        for (let tran of transactions) {
            vendorOpts.push(tran.vendor);
        }

        setVendorOptions(vendorOpts);
    }, [transactions]);

    let vendorField = register("vendor", {
        required: true,
        maxLength: 80,
    });

    const vendorChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        vendorField.onChange(e);
        console.log("customthing " + e.target.value);
        let transactionWithVendor = transactions?.find(
            (x) => x.vendor == e.target.value
        );

        console.log(transactionWithVendor);

        if (transactionWithVendor?.category_id !== undefined) {
            setValue("category", transactionWithVendor.category_id);
        }
    };

    return (
        <div className=" max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div className="grid gap-2">
                    <datalist id="vendorInputList">
                        {vendorOptions?.map((x) => {
                            return <option key={x} value={x}></option>;
                        })}
                    </datalist>
                    <label htmlFor="vendor">Vendor</label>
                    <input
                        id="vendor"
                        autoComplete="off"
                        type="text"
                        placeholder="Walmart"
                        {...vendorField}
                        onChange={(e) => {
                            vendorChange(e);
                        }}
                        list="vendorInputList"
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

                    <label htmlFor="transactionType">Transaction Type</label>
                    <select
                        id="transactionType"
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        {...register("type")}
                    >
                        {[
                            { name: "income", label: "Income" },
                            { name: "expense", label: "Expense" },
                        ]?.map((c) => {
                            return (
                                <option key={c.name} value={c.name}>
                                    {c.label}
                                </option>
                            );
                        })}
                    </select>
                </div>
                <br />
                <input
                    type="submit"
                    className="float-right bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md cursor:pointer w-36"
                />
            </form>
        </div>
    );
};
