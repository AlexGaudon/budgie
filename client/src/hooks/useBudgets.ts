import { useQuery, useMutation, useQueryClient } from "react-query";

import { z } from "zod";

import { type Budget, budgetSchema } from "../types";
import { CreateBudgetForm } from "../components/AddBudget";

const fetchBudgets = async () => {
    const res = await fetch("/api/budgets");

    if (res.ok) {
        const data = await res.json();
        if ("data" in data) {
            const parsed = z.array(budgetSchema).parse(data.data);
            return parsed;
        }
    } else {
        throw new Error("Error fetching budgets");
    }
};

export const useCreateBudgetMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (newBudget: CreateBudgetForm) => {
        const res = await fetch("/api/budgets", {
            method: "POST",
            body: JSON.stringify({
                category: newBudget.category,
                period: newBudget.period,
                amount: Number(
                    Math.round(
                        parseFloat(newBudget.amount.toString()) * 100
                    ).toString()
                ),
            }),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            queryClient.invalidateQueries("budgets");
        } else {
            throw new Error("Error creating budget");
        }
    });
};

export const useUpdateBudgetMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (updatedBudget: Budget) => {
        const res = await fetch(`/api/budgets/${updatedBudget.id}`, {
            method: "PUT",
            body: JSON.stringify(updatedBudget),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            queryClient.invalidateQueries("budgets");
        } else {
            throw new Error("Error updating budget");
        }
    });
};

export const useDeleteBudgetMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (budgetId: string) => {
        const res = await fetch(`/api/budgets/${budgetId}`, {
            method: "DELETE",
        });

        if (res.ok) {
            queryClient.invalidateQueries("budgets");
        } else {
            throw new Error("Error deleting budget");
        }
    });
};

export const useBudgetsQuery = () => {
    return useQuery("budgets", fetchBudgets);
};
