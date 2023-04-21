import { useContext, useState, useEffect } from "react";
import { AuthContext } from "../contexts/AuthContext";

export const useAuth = () => {
    const authContext = useContext(AuthContext);
    const [isLoading, setIsLoading] = useState(true);

    if (!authContext) {
        throw new Error("useAuth must be used within an AuthProvider");
    }

    const { isLoggedIn, user, login, logout } = authContext;

    return {
        isLoading,
        isLoggedIn,
        user,
        login,
        logout,
    };
};
