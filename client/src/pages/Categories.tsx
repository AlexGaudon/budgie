import { useEffect } from "react";

import { useCategoriesQuery } from "../hooks/useCategories";

type CategoryProps = {
    name: string;
};

export const Category = ({ name }: CategoryProps) => {
    return <h1>{name}</h1>;
};

export const Categories = () => {
    const { data: categories, isLoading, error } = useCategoriesQuery();

    if (isLoading) {
        return <h1>Loading...</h1>;
    }

    if (error) {
        console.log(error);
        return <h1>Error.</h1>;
    }

    return (
        <div className="container block">
            {categories?.map((c) => {
                return <Category key={c.id} name={c.name}></Category>;
            })}
        </div>
    );
};
