import React, { useEffect, useState } from "react";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Navigate, useNavigate } from "react-router-dom";

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
        return <Navigate to="/" />;
    }

    const onSubmit: SubmitHandler<LoginForm> = async (data) => {
        auth.login(data.username, data.password);
    };

    return (
        <div>
            <form onSubmit={handleSubmit(onSubmit)}>
                <input
                    type="text"
                    placeholder="Username"
                    {...register("username", { required: true, maxLength: 80 })}
                />
                {errors.username && <p>User name is required.</p>}

                <input
                    type="password"
                    placeholder="Password"
                    {...register("password", {
                        required: true,
                        maxLength: 100,
                    })}
                />
                {errors.password && <p>Password is required.</p>}

                <input type="submit" />
            </form>
        </div>
    );
};
