import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { optional, z } from "zod";
import { useCreateBudgetMutation } from "../hooks/useBudgets";
import { useCategoriesQuery } from "../hooks/useCategories";

const createBudgetSchema = z.object({
    category: z.string(),
    amount: z.string(),
    period: z.string(),
});

export type CreateBudgetForm = z.infer<typeof createBudgetSchema>;

export const AddBudget = ({ onFinish }: { onFinish: () => void }) => {
    const {
        register,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm<CreateBudgetForm>({
        resolver: zodResolver(createBudgetSchema),
    });

    const createBudgetMutation = useCreateBudgetMutation();

    const {
        data: categories,
        isLoading: categoryLoading,
        error: categoryError,
    } = useCategoriesQuery();

    const onSubmit: SubmitHandler<CreateBudgetForm> = async (data) => {
        data.period = new Date(data.period).toISOString();
        let res = await createBudgetMutation.mutateAsync(data);

        console.log(res);

        onFinish();
        reset();
    };

    return (
        <div className="flex items-center justify-center h-full w-8/12">
            <form onSubmit={handleSubmit(onSubmit)}>
                <div className="grid grid-cols-2 gap-4">
                    <label htmlFor="category">Category</label>
                    {errors.category && "Error: " + errors.category}
                    <select
                        id="category"
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        {...register("category", {
                            required: true,
                            maxLength: 80,
                        })}
                    >
                        {!categoryLoading &&
                            categories?.map((c) => {
                                return (
                                    <option key={c.id} value={c.id}>
                                        {c.name}
                                    </option>
                                );
                            })}
                        {categoryLoading && <option>Loading...</option>}
                    </select>

                    <label htmlFor="amount">Amount</label>
                    <input
                        id="amount"
                        autoComplete="off"
                        type="number"
                        step="0.01"
                        placeholder="0.00"
                        {...register("amount", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />

                    <label htmlFor="period">Budget Period</label>
                    <input
                        id="period"
                        autoComplete="off"
                        type="month"
                        placeholder="Date"
                        {...register("period", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />

                    <input
                        type="submit"
                        className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md cursor:pointer"
                    />
                </div>
            </form>
        </div>
    );
};
