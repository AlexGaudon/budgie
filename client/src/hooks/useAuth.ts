import { useContext, useState, useEffect } from "react";
import { AuthContext } from "../contexts/AuthContext";

export const useAuth = () => {
    const authContext = useContext(AuthContext);

    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }

    const { isLoggedIn, user, login, logout, isLoading } = authContext;

    return {
        isLoading,
        isLoggedIn,
        user,
        login,
        logout,
    };
};
