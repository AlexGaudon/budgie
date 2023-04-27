import { useForm } from "react-hook-form";
import { useCategoriesQuery } from "../hooks/useCategories";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { AddCategory } from "../components/AddCategory";

type CategoryProps = {
    id: string;
    name: string;
};

export const Category = ({ id, name }: CategoryProps) => {
    return (
        <div
            key={id}
            className="border border-gray-300 rounded-md p-4 mb-4 flex justify-between items-center"
        >
            <div>
                <p className="text-gray-600 text-sm">ID: {id}</p>
                <p>{name}</p>
            </div>
            <button
                onClick={() => {}}
                className="px-4 py-2 bg-red-500 text-white rounded-md"
            >
                Delete
            </button>
        </div>
    );
};

export const Categories = () => {
    const { data: categories, isLoading, error } = useCategoriesQuery();
    const [isCreating, setIsCreating] = useState(false);
    if (isLoading) {
        return <h1>Loading...</h1>;
    }

    if (error) {
        console.log(error);
        return <h1>Error.</h1>;
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
