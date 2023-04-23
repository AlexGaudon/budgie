import React, { createContext, useState, useEffect } from "react";

import { z } from "zod";

const userSchema = z.object({
    userId: z.string(),
    username: z.string(),
});

type User = z.TypeOf<typeof userSchema>;

const errorSchema = z.object({
    message: z.string(),
});

interface AuthContextType {
    isLoggedIn: boolean;
    isLoading: boolean;
    user: User | null;
    login: (username: string, password: string) => void;
    logout: () => void;
}

export const AuthContext = createContext<AuthContextType>({
    isLoggedIn: false,
    isLoading: true,
    user: null,
    login: () => {},
    logout: () => {},
});

export const AuthProvider = ({
    children,
}: {
    children: React.ReactElement | React.ReactElement[];
}) => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState("");

    const login = async (username: string, password: string) => {
        setIsLoading(true);
        let res = await fetch("/api/user/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username,
                password,
            }),
        });

        if (res.ok) {
            let validatedUser = userSchema.parse(await res.json());
            setUser(validatedUser);
            setIsLoggedIn(true);
            setIsLoading(false);
        } else {
            let validatedError = errorSchema.parse(await res.json());
            setError(validatedError.message);
            setIsLoggedIn(false);
            setIsLoading(false);
            setUser(null);
        }
    };

    const logout = async () => {
        setIsLoading(true);
        let res = await fetch("/api/user/logout");
        if (res.ok) {
            setIsLoggedIn(false);
            setUser(null);
        }
        setError("");
        setIsLoading(false);
    };

    useEffect(() => {
        const checkAuthenticationStatus = async () => {
            setIsLoading(true);
            let res = await fetch("/api/user/me");
            if (res.ok) {
                let body = await res.json();
                console.log(body);
                const validatedUser = userSchema.parse(body);
                setUser(validatedUser);
                setIsLoggedIn(true);
            } else {
                console.log(await res.json());
                setUser(null);
                setIsLoggedIn(false);
            }
            setIsLoading(false);
        };

        checkAuthenticationStatus();
    }, []);

    return (
        <AuthContext.Provider
            value={{ isLoggedIn, user, login, logout, isLoading }}
        >
            {children}
        </AuthContext.Provider>
    );
};
