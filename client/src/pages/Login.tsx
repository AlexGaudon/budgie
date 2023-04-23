import React from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Navigate } from "react-router-dom";
import { z } from "zod";
import { useAuth } from "../hooks/useAuth";

const loginFormSchema = z.object({
    username: z.string(),
    password: z.string(),
});

type LoginForm = z.infer<typeof loginFormSchema>;

export const Login = () => {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginForm>({ resolver: zodResolver(loginFormSchema) });

    const auth = useAuth();

    if (auth.isLoggedIn && !auth.isLoading) {
        return <Navigate to="/"></Navigate>;
    }

    const onSubmit: SubmitHandler<LoginForm> = async (data) => {
        auth.login(data.username, data.password);
    };

    return (
        <div className="flex items-center justify-center h-full">
            <div className="bg-white p-8 rounded-md shadow-md w-full max-w-sm">
                <h1 className="text-3xl font-bold mb-6 text-center text-black">
                    Login
                </h1>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <input
                        type="text"
                        placeholder="Username"
                        {...register("username", {
                            required: true,
                            maxLength: 80,
                        })}
                        className="w-full px-4 py-3 mb-3 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    {errors.username && (
                        <p className="text-red-500 mb-2">
                            Username is required.
                        </p>
                    )}

                    <input
                        type="password"
                        placeholder="Password"
                        {...register("password", {
                            required: true,
                            maxLength: 100,
                        })}
                        className="w-full px-4 py-3 mb-3 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                    {errors.password && (
                        <p className="text-red-500 mb-2">
                            Password is required.
                        </p>
                    )}

                    <button
                        type="submit"
                        className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-6 rounded-md w-full"
                    >
                        Log in
                    </button>
                </form>
            </div>
        </div>
    );
};
