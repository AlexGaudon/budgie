import { useEffect, useState } from "react";
import { type Budget, budgetSchema } from "../types";

import { z } from "zod";
import { Navigate } from "react-router-dom";

import { useAuth } from "../hooks/useAuth";

const budgetResponse = z.object({
    data: z.array(budgetSchema),
});

export const Budgets = () => {
    const [budgets, setBudgets] = useState<Budget[]>([]);
    const [refresh, setRefresh] = useState(false);

    let auth = useAuth();

    if (!auth.isLoggedIn && !auth.isLoading) {
        return <Navigate to="/login" />;
    }

    useEffect(() => {
        async function getBudgets() {
            try {
                const res = await fetch("/api/budgets", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });
                if (res.ok) {
                    console.log("was ok");
                    const data = await res.json();
                    const parsedData = budgetResponse.parse(data).data;
                    setBudgets(parsedData);

                    console.log(JSON.stringify(parsedData));
                } else {
                    console.log("oh fk");
                }
            } catch (error) {
                console.log(error);
            }
        }

        getBudgets();
    }, [refresh]);

    const deleteBudget = async (id: string) => {
        let res = await fetch(`http://localhost:3000/api/budgets/${id}`, {
            method: "DELETE",
        });
        if (res.ok) {
            setRefresh((v) => !v);
        }
    };

    return (
        <div>
            {budgets.map((b) => {
                return (
                    <div key={b.id}>
                        <p>({b.id})</p>
                        <p>
                            {b.category}(
                            {b.recurring ? "recurring" : "not recurring"}) for{" "}
                            {b.amount}
                        </p>
                        <button
                            onClick={() => {
                                deleteBudget(b.id);
                            }}
                        >
                            delete
                        </button>
                        <br />
                    </div>
                );
            })}
        </div>
    );
};
