import {
    useCategoriesQuery,
    useDeleteCategoryMutation,
} from "../hooks/useCategories";
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

    const showCategory = () => {
        console.log("CALLING");
        return <Navigate to={`transactions`} />;
    };
    return (
        <div
            key={id}
            className="border border-gray-300 rounded-md p-4 mb-4 flex justify-between items-center max-w-md"
        >
            <div>
                <p className="text-gray-600 text-sm">ID: {id}</p>
                <p>{name}</p>
            </div>
            <button
                onClick={() => {
                    showCategory();
                }}
            >
                VIEW
            </button>
            <button
                onClick={() => {
                    deleteCategory.mutateAsync(id);
                }}
                className="px-4 py-2 bg-red-500 text-white rounded-md"
            >
                Delete
            </button>
        </div>
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

            {categories?.map((c) => {
                return <Category key={c.id} id={c.id} name={c.name}></Category>;
            })}
        </div>
    );
};
