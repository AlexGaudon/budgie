import React, { MouseEventHandler, useEffect, useState } from "react";

import { Transaction } from "../types";

import { useCategoriesQuery } from "../hooks/useCategories";

type RowProps = {
    index: number;
    editing: boolean;
    row: Transaction;
    onRowEvent: (eventType: "edit" | "delete", subject: Transaction) => void;
    setEditingIndex: (i: number) => void;
};

const normalizeDate = (date: Date | string) => {
    if (typeof date == "string") {
        return date;
    }
    return date.toISOString().slice(0, 10);
};

export const TransactionRow = ({
    index,
    row,
    onRowEvent,
    setEditingIndex,
    editing,
}: RowProps) => {
    const { data: categories, isLoading, error } = useCategoriesQuery();
    const [editedRow, setEditedRow] = useState(row);
    const getClasses = () => {
        return "border border-gray-200";
    };

    const handleChange = (
        event: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
        field:
            | "date"
            | "category_id"
            | "description"
            | "amount"
            | "type"
            | "vendor"
    ) => {
        setEditedRow((prevEditedRow) => {
            const er = { ...prevEditedRow };
            if (field === "date") {
                er.date = new Date(event.target.value);
            } else {
                er[field] = event.target.value;
            }
            return er;
        });
    };

    const onBlur = () => {
        onRowEvent("edit", editedRow);
    };

    const handleRowClick = (
        event: React.MouseEvent<HTMLTableRowElement, MouseEvent>
    ) => {
        if ((event.target as HTMLElement).tagName.toLowerCase() !== "button") {
            setEditingIndex(index);
        }
    };

    const handleDelete = () => {
        onRowEvent("delete", row);
    };

    return (
        <tr className={getClasses()} onClick={handleRowClick}>
            <td className="px-4 py-2 text-left">
                {editing ? (
                    <input
                        type="date"
                        value={normalizeDate(editedRow.date)}
                        onChange={(e) => {
                            handleChange(e, "date");
                        }}
                        onBlur={onBlur}
                    />
                ) : (
                    normalizeDate(editedRow.date)
                )}
            </td>
            <td className="px-4 py-2 text-left">
                {editing ? (
                    <input
                        type="text"
                        value={editedRow.vendor}
                        onChange={(e) => {
                            handleChange(e, "vendor");
                        }}
                        onBlur={onBlur}
                    />
                ) : (
                    editedRow.vendor
                )}
            </td>
            <td className="px-4 py-2 text-left">
                {editing ? (
                    <input
                        type="text"
                        value={
                            editedRow.description != null
                                ? editedRow.description
                                : ""
                        }
                        onChange={(e) => {
                            handleChange(e, "description");
                        }}
                        onBlur={onBlur}
                    />
                ) : (
                    editedRow.description
                )}
            </td>
            <td className="px-4 py-2 text-left">
                {editing ? (
                    <select
                        name="category"
                        id="category"
                        onChange={(e) => {
                            handleChange(e, "category_id");
                        }}
                        onBlur={onBlur}
                        value={editedRow.category_id}
                    >
                        {categories?.map((c) => {
                            return (
                                <option key={c.id} value={c.id}>
                                    {c.name}
                                </option>
                            );
                        })}
                    </select>
                ) : (
                    editedRow.category
                )}
            </td>
            <td
                className={
                    editedRow.type === "income"
                        ? "text-green-400"
                        : "" && "px-4 py-2 text-left"
                }
            >
                {editing ? (
                    <input
                        type="number"
                        value={editedRow.amount}
                        onChange={(e) => handleChange(e, "amount")}
                        onBlur={onBlur}
                    />
                ) : (
                    (editedRow.type === "income" && "$" + editedRow.amount) ||
                    "-$" + editedRow.amount
                )}
            </td>
            <td className="px-4 py-2 text-left">
                <button onClick={handleDelete}>DELETE</button>
            </td>
        </tr>
    );
};
