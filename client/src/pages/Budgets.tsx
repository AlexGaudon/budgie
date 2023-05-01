import { useEffect, useState } from "react";
import { budgetSchema } from "../types";

import { z } from "zod";
import { Navigate } from "react-router-dom";

import { useAuth } from "../hooks/useAuth";
import { useBudgetsQuery, useDeleteBudgetMutation } from "../hooks/useBudgets";
import { AddBudget } from "../components/AddBudget";

type BudgetProps = {
    id: string;
    name: string;
    period: string;
    amount: string;
};

export const Budget = ({ id, name, period, amount }: BudgetProps) => {
    const deleteBudget = useDeleteBudgetMutation();
    return (
        <div
            key={id}
            className="border border-gray-300 rounded-md p-4 mb-4 flex justify-between items-center"
        >
            <div>
                <p className="text-gray-600 text-sm">ID: {id}</p>
                <p>
                    {name} for {amount}
                </p>
                <p>for the month of {period}</p>
            </div>
            <button
                onClick={() => {
                    deleteBudget.mutateAsync(id);
                }}
                className="px-4 py-2 bg-red-500 text-white rounded-md"
            >
                Delete
            </button>
        </div>
    );
};

export const Budgets = () => {
    let auth = useAuth();

    const { data: budgets, isLoading, error } = useBudgetsQuery();
    const [isCreating, setIsCreating] = useState(false);
    if (isLoading) {
        return <h1>Loading...</h1>;
    }

    if (error) {
        console.log(error);
        return <h1>Error: {"" + error}</h1>;
    }

    if (!auth.isLoggedIn && !auth.isLoading) {
        return <Navigate to="/login" />;
    }

    return (
        <div>
            <button
                onClick={() => {
                    setIsCreating(true);
                }}
            >
                Add Budget
            </button>
            {isCreating && (
                <AddBudget
                    onFinish={() => {
                        setIsCreating(false);
                    }}
                ></AddBudget>
            )}

            {budgets?.map((budget) => {
                return (
                    <Budget
                        key={budget.id}
                        id={budget.id}
                        name={budget.name}
                        amount={budget.amount}
                        period={budget.period}
                    ></Budget>
                );
            })}
        </div>
    );
};
