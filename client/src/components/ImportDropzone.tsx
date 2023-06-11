import { useState } from "react";
import { useDropzone } from "react-dropzone";
import { Category } from "../types";
import Papa from "papaparse";
import { useCreateTransactionMutation } from "../hooks/useTransactions";
import { useCategoriesQuery } from "../hooks/useCategories";

type csvTransaction = {
    date: string;
    transactionDetails: string;
    fundsOut: string;
    fundsIn: string;
    creditCard: string;
};

export const ImportDropZone = () => {
    const [files, setFiles] = useState<File[]>([]);

    const { data: categories, isLoading, isError } = useCategoriesQuery();

    const createTransaction = useCreateTransactionMutation();

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop: (acceptedFiles) => {
            setFiles(acceptedFiles);
        },
    });
    const handleProcessFiles = () => {
        files.forEach((file) => {
            const reader = new FileReader();
            reader.onload = () => {
                let category = categories?.find(
                    (e: Category) => e.name === "Uncategorized"
                );

                let id = category?.id;
                const result = Papa.parse<csvTransaction>(
                    reader.result as string,
                    {
                        header: true,
                        transformHeader: (header: string) =>
                            header
                                .trim()
                                .split(" ")
                                .map((word, index) =>
                                    index === 0
                                        ? word.toLowerCase()
                                        : `${word.charAt(0).toUpperCase()}${word
                                              .slice(1)
                                              .toLowerCase()}`
                                )
                                .join(""),
                    }
                );
                if (id === undefined) return;
                result.data.forEach((d) => {
                    console.log(d);
                    if (d.fundsOut === "") return;
                    createTransaction.mutateAsync({
                        vendor: d.transactionDetails,
                        amount: d.fundsOut,
                        date: d.date,
                        type: "expense",
                        description: "",
                        category: id || "",
                    });
                });
                console.log(result.data);
            };
            reader.readAsText(file);
        });
    };

    return (
        <>
            <div
                {...getRootProps()}
                className={`border-2 border-dashed p-4 ${
                    isDragActive ? "border-gray-300" : "border-blue-300"
                }`}
            >
                <input {...getInputProps()} />
                <p>Drop Transactions Here</p>
                {files.length > 0 && (
                    <ul>
                        {files.map((file) => (
                            <li key={file.name}>{file.name}</li>
                        ))}
                    </ul>
                )}
            </div>
            <button
                onClick={handleProcessFiles}
                disabled={files.length === 0}
                className={`px-4 py-2 rounded-md ${
                    files.length === 0
                        ? "bg-gray-300 cursor-not-allowed"
                        : "bg-blue-500 hover:bg-blue-600 text-white"
                }`}
            >
                Process Files
            </button>
        </>
    );
};
