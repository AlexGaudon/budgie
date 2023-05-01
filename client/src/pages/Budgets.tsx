import { useEffect, useState } from "react";
import { budgetSchema } from "../types";

import { z } from "zod";
import { Navigate } from "react-router-dom";

import { useAuth } from "../hooks/useAuth";
import { useBudgetsQuery, useDeleteBudgetMutation } from "../hooks/useBudgets";
import { AddBudget } from "../components/AddBudget";

type BudgetProps = {
    id: string;
    category: string;
    period: string;
    amount: string;
    utilization: string;
};

export const Budget = ({
    id,
    category,
    period,
    amount,
    utilization,
}: BudgetProps) => {
    const deleteBudget = useDeleteBudgetMutation();

    const amountLeft = () => {
        // todo
        return "$100 left";
    };
    return (
        <div
            key={id}
            className="border border-gray-300 rounded-md p-4 mb-4 flex items-center max-w-md"
        >
            <div className="w-5/6">
                <p>
                    {category}
                    <span className="float-right">{amountLeft()}</span>
                </p>
                <p>
                    <span>{utilization}</span> of{" "}
                    <span className="text-amber-300">{amount}</span>
                </p>
            </div>
            <button
                onClick={() => {
                    deleteBudget.mutateAsync(id);
                }}
                className="ml-2 px-4 py-2 bg-red-500 text-white rounded-md"
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
        return <h1>{"" + error}</h1>;
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
                        category={budget.category}
                        amount={budget.amount}
                        period={budget.period}
                        utilization="$TODO"
                    ></Budget>
                );
            })}
        </div>
    );
};
