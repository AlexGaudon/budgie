import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";
import {
    useCategoriesQuery,
    useCreateCategoryMutation,
} from "../hooks/useCategories";

const categoryCreateSchema = z.object({
    name: z.string(),
});

export type CreateCategoryForm = z.infer<typeof categoryCreateSchema>;

export const AddCategory = ({ onFinish }: { onFinish: () => void }) => {
    const {
        register,
        handleSubmit,
        formState: { errors },
        reset,
    } = useForm<CreateCategoryForm>({
        resolver: zodResolver(categoryCreateSchema),
    });

    const {
        data: categories,
        isLoading: isLoading,
        error: error,
    } = useCategoriesQuery();

    const createCategoryMutation = useCreateCategoryMutation();

    const onSubmit: SubmitHandler<CreateCategoryForm> = async (data) => {
        console.log(data);
        let input = {
            name: data.name,
        };
        await createCategoryMutation.mutateAsync(input);

        onFinish();
    };

    return (
        <div className="flex items-center justify-center h-full w-8/12">
            {errors && <h1>{JSON.stringify(errors)}</h1>}
            <form onSubmit={handleSubmit(onSubmit)}>
                <div className="grid grid-cols-2 gap-4">
                    <label htmlFor="name">Name</label>
                    <input
                        id="name"
                        autoComplete="off"
                        placeholder="My Category"
                        type="text"
                        {...register("name", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="border border-gray-300 rounder-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
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
