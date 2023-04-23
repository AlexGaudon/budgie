import { z } from "zod";

const transactionSchema = z.object({
    id: z.string(),
    userid: z.string(),
    description: z.string(),
    category: z.string(),
    amount: z.number().transform((num) => {
        // Check if the input number is valid
        if (isNaN(num)) {
            throw new Error("Invalid amount value");
        }

        num /= 100;
        return num.toFixed(2);
    }),
    type: z.string().refine(
        (str) => {
            return str == "income" || str == "expense";
        },
        {
            message: 'Type must be equal to "income" or "expense"',
        }
    ),
    date: z.string().transform((str) => new Date(str)),
    updated_at: z.string(),
    created_at: z.string(),
    deleted_at: z.string().nullable(),
});

type Transaction = z.infer<typeof transactionSchema>;

const newTransactionSchema = z.object({
    description: z.string(),
    category: z.string(),
    amount: z.number(),
    type: z.string().refine(
        (str) => {
            return str == "income" || str == "expense";
        },
        {
            message: 'Type must be equal to "income" or "expense"',
        }
    ),
    date: z.string(),
});

type NewTransaction = z.infer<typeof newTransactionSchema>;

const budgetSchema = z.object({
    id: z.string(),
    userid: z.string(),
    category: z.string(),
    amount: z.number().transform((num) => {
        // Check if the input number is valid
        if (isNaN(num)) {
            throw new Error("Invalid amount value");
        }

        num /= 100;

        // Convert the number to currency format
        return num.toLocaleString("en-US", {
            style: "currency",
            currency: "USD",
        });
    }),
    recurring: z.boolean(),
    updated_at: z.string(),
    created_at: z.string(),
    deleted_at: z.string().nullable(),
});

type Budget = z.infer<typeof budgetSchema>;

const userSchema = z.object({
    UserID: z.string(),
    Username: z.string(),
});

type User = z.infer<typeof userSchema>;

const authErrorSchema = z.object({
    message: z.string(),
});

type AuthError = z.infer<typeof authErrorSchema>;

export {
    type Transaction,
    transactionSchema,
    type Budget,
    budgetSchema,
    type NewTransaction,
    newTransactionSchema,
};