import React, { MouseEventHandler, useState } from "react";

import { Transaction } from "../types";

import { z } from "zod";

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
    const [editedRow, setEditedRow] = useState(row);
    const getClasses = () => {
        return "border border-gray-200";
    };

    const handleChange = (
        event: React.ChangeEvent<HTMLInputElement>,
        field: "date" | "category" | "description" | "amount" | "type"
    ) => {
        setEditedRow((prevEditedRow) => {
            const er = { ...prevEditedRow }; // Create a shallow copy of the state object
            if (field === "date") {
                er.date = new Date(event.target.value);
            } else {
                er[field] = event.target.value;
            }
            return er; // Return the updated state object
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
            <td>{editedRow.id.substring(0, 4)}</td>
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
                        value={editedRow.description}
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
                    <input
                        type="text"
                        value={editedRow.category}
                        onChange={(e) => {
                            handleChange(e, "category");
                        }}
                        onBlur={onBlur}
                    />
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
