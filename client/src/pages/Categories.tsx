import { useForm } from "react-hook-form";
import {
    useCategoriesQuery,
    useDeleteCategoryMutation,
} from "../hooks/useCategories";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { AddCategory } from "../components/AddCategory";
import { useAuth } from "../hooks/useAuth";
import { Navigate } from "react-router-dom";

type CategoryProps = {
    id: string;
    name: string;
};

export const Category = ({ id, name }: CategoryProps) => {
    const deleteCategory = useDeleteCategoryMutation();

    return (
        <tr key={id}>
            <td className="border px-4 py-2">{name}</td>
            <td className="border px-4 py-2">
                <button
                    className="px-4 py-2 bg-red-500 text-white rounded-md"
                    onClick={() => deleteCategory.mutateAsync(id)}
                >
                    Delete
                </button>
            </td>
        </tr>
    );
};

export const Categories = () => {
    let auth = useAuth();
    const { data: categories, isLoading, error } = useCategoriesQuery();
    const [isCreating, setIsCreating] = useState(false);
    if (isLoading) {
        return <h1>Loading...</h1>;
    }

    if (error) {
        console.log(error);
        return <h1>Error.</h1>;
    }

    if (!auth.isLoggedIn && !auth.isLoading) {
        return <Navigate to="/login" />;
    }

    return (
        <div>
            <button
                className="px-4 py-4 rounded-md bg-blue-500 hover:bg-blue-600"
                onClick={() => {
                    setIsCreating(true);
                }}
            >
                Add New Category
            </button>
            {isCreating && (
                <AddCategory
                    onFinish={() => {
                        setIsCreating(false);
                    }}
                ></AddCategory>
            )}

            <table className="table-auto w-auto">
                <thead>
                    <tr className="border">
                        <th className="px-4 py-2">Category Name</th>
                        <th className="px-4 py-2">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {categories?.map((c) => {
                        return (
                            <Category
                                key={c.id}
                                id={c.id}
                                name={c.name}
                            ></Category>
                        );
                    })}
                </tbody>
            </table>
        </div>
    );
};
