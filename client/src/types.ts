import { TypeOf, z } from "zod";

const transactionSchema = z.object({
    id: z.string(),
    userid: z.string(),
    vendor: z.string(),
    description: z.string(),
    category_id: z.string(),
    category_name: z.string(),
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
    vendor: z.string(),
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

const categorySchema = z.object({
    id: z.string(),
    userid: z.string(),
    name: z.string(),
});

type Category = z.infer<typeof categorySchema>;

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

export {
    type Transaction,
    transactionSchema,
    type Budget,
    budgetSchema,
    type NewTransaction,
    newTransactionSchema,
    type Category,
    categorySchema,
};
