import { useState } from "react";

import { Link, Navigate } from "react-router-dom";

import { useAuth } from "../hooks/useAuth";
import { useBudgetsUtilizationQuery } from "../hooks/useBudgets";
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
    const getAmtAsNums = () => {
        let amtNum = parseFloat(amount.replace("$", "").replace(",", ""));
        let utilNum = parseFloat(utilization.replace("$", "").replace(",", ""));

        return { amtNum, utilNum };
    };
    const amountRemaining = () => {
        let { amtNum, utilNum } = getAmtAsNums();
        return `$${(amtNum - utilNum).toLocaleString("en-US", {
            currency: "USD",
        })}`;
    };

    const utilPercent = () => {
        let { amtNum, utilNum } = getAmtAsNums();
        return (utilNum / amtNum) * 100;
    };

    return (
        <div className="border border-gray-300 p-4 mb-4 flex items-center max-w-md">
            <div className="w-full">
                <p>
                    <Link
                        to={`/transactions?filter=${encodeURIComponent(
                            category
                        )}`}
                    >
                        {category}
                    </Link>
                    <span className="float-right">
                        {amountRemaining()} left
                    </span>
                </p>

                <div className="w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700">
                    <div
                        className={
                            utilPercent() < 100
                                ? "bg-blue-300 h-2.5 rounded-full"
                                : "bg-red-400 h-2.5 rounded-full"
                        }
                        style={{
                            width:
                                utilPercent() < 100
                                    ? utilPercent() + "%"
                                    : 100 + "%",
                        }}
                    ></div>
                </div>

                <p>
                    {utilization} of {amount}
                </p>
            </div>
        </div>
    );
};

export const Budgets = () => {
    let auth = useAuth();

    const [year, setYear] = useState(new Date().getFullYear().toString());
    const [month, setMonth] = useState(
        new Date().getMonth() < 10
            ? "0" + (new Date().getMonth() + 1)
            : new Date().getMonth() + 1
    );

    const getBudgets = useBudgetsUtilizationQuery();
    const { data: budgets, isLoading, error } = getBudgets(year + "-" + month);

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
            <div className="flex items-center space-x-4">
                <select
                    className="rounded-md shadow-sm border-gray-300 focus:border-indigo-500 focus:ring-indigo-500"
                    value={year}
                    onChange={(e) => {
                        setYear(e.target.value);
                    }}
                >
                    <option value="">Year</option>
                    <option value="2022">2022</option>
                    <option value="2023">2023</option>
                    <option value="2024">2024</option>
                </select>

                <select
                    className="rounded-md shadow-sm border-gray-300 focus:border-indigo-500 focus:ring-indigo-500"
                    value={month}
                    onChange={(e) => {
                        setMonth(e.target.value);
                    }}
                >
                    <option value="">Month</option>
                    <option value="01">January</option>
                    <option value="02">February</option>
                    <option value="03">March</option>
                    <option value="04">April</option>
                    <option value="05">May</option>
                    <option value="06">June</option>
                    <option value="07">July</option>
                    <option value="08">August</option>
                    <option value="09">September</option>
                    <option value="10">October</option>
                    <option value="11">November</option>
                    <option value="12">December</option>
                </select>
            </div>

            <button
                className="bg-blue-300 text-gray-900 px-3 py-2 rounded-md text-sm font-medium mb-4"
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
            <div className="grid md:grid-cols-3 gap-4">
                {budgets?.map((budget) => {
                    return (
                        <Budget
                            key={budget.id}
                            id={budget.id}
                            category={budget.category}
                            amount={budget.amount}
                            period={budget.period}
                            utilization={budget.utilization}
                        ></Budget>
                    );
                })}
            </div>
        </div>
    );
};
