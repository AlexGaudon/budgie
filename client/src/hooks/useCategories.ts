import { useQuery, useMutation, useQueryClient } from "react-query";

import { z } from "zod";

import { type Category, categorySchema } from "../types";

const fetchCategories = async () => {
    const res = await fetch("/api/categories");

    if (res.ok) {
        const data = await res.json();
        if ("data" in data) {
            const parsed = z.array(categorySchema).parse(data.data);
            return parsed;
        }
    } else {
        throw new Error("Error fetching categories");
    }
};

export const useCreateCategoryMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (newCategory: { name: string }) => {
        const res = await fetch("/api/categories", {
            method: "POST",
            body: JSON.stringify(newCategory),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            queryClient.invalidateQueries("categories");
        } else {
            throw new Error("Error creating category");
        }
    });
};

export const useUpdateCategoryMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (updatedCategory: Category) => {
        const res = await fetch(`/api/categories/${updatedCategory.id}`, {
            method: "PUT",
            body: JSON.stringify(updatedCategory),
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.ok) {
            queryClient.invalidateQueries("categories");
        } else {
            throw new Error("Error updating category");
        }
    });
};

export const useDeleteCategoryMutation = () => {
    const queryClient = useQueryClient();

    return useMutation(async (categoryId: string) => {
        const res = await fetch(`/api/categories/${categoryId}`, {
            method: "DELETE",
        });

        if (res.ok) {
            queryClient.invalidateQueries("categories");
        } else {
            throw new Error("Error deleting category");
        }
    });
};

export const useCategoriesQuery = () => {
    return useQuery("categories", fetchCategories);
};
