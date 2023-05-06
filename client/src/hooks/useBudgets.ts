import { useQuery, useMutation, useQueryClient } from "react-query";

import { z } from "zod";

import { type Budget, budgetSchema } from "../types";
import { CreateBudgetForm } from "../components/AddBudget";

const budgetUtilSchema = z.object({
    id: z.string(),
    created_at: z.string(),
    updated_at: z.string(),
    user: z.string(),
    category: z.string(),
    amount: z.number().transform((num) => {
        // Check if the input number is valid
        if (isNaN(num)) {
            throw new Error("Invalid amount value");
        }

        num /= 100;

        // Convert the number to currency format
        return num.toLocaleString("en-US", {
            style: "currency",
            currency: "USD",
        });
    }),
    period: z.string().transform((input) => {
        return new Date(input).toISOString().substring(0, 7);
    }),
    utilization: z.number().transform((num) => {
        // Check if the input number is valid
        if (isNaN(num)) {
            throw new Error("Invalid amount value");
        }

        num /= 100;

        // Convert the number to currency format
        return num.toLocaleString("en-US", {
            style: "currency",
            currency: "USD",
        });
    }),
});

const fetchBudgetsByPeriod = async (period: string) => {
    const res = await fetch(`/api/budgets/utilization/${period}`);

    if (res.ok) {
        const data = await res.json();
        if ("data" in data) {
            const parsed = z.array(budgetUtilSchema).parse(data.data);
            return parsed;
        }
    } else {
        throw new Error("Error fetching budgets");
    }
};

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

export const useBudgetsUtilizationQuery = () => {
    return (period: string) => {
        return useQuery(["budgetsItl", period], () =>
            fetchBudgetsByPeriod(period)
        );
    };
};
