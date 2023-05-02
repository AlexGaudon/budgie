import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

interface NavLinkProps {
    to: string;
    exact?: boolean;
    children: React.ReactNode;
}

const NavLink = ({ to, exact = false, children }: NavLinkProps) => {
    let className =
        "text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium";
    let activeClassName =
        "bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium";
    const location = useLocation();
    const isActive = exact
        ? location.pathname === to
        : location.pathname.startsWith(to);

    return (
        <Link to={to} className={isActive ? activeClassName : className}>
            {children}
        </Link>
    );
};

export const NavBar = () => {
    let auth = useAuth();
    const [isMobileMenuOpen, setMobileMenuOpen] = useState(false);

    const toggleMobileMenu = () => {
        setMobileMenuOpen(!isMobileMenuOpen);
    };

    return (
        <nav className="bg-gray-800 stick top-0 w-full border-b-2 border-slate-800 shadow-md">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center">
                        <div className="flex-shrink-0">
                            <Link to="/">
                                Bu<span className="text-blue-500">dg</span>ie
                            </Link>
                        </div>
                        <div className="hidden md:block">
                            <div className="ml-10 flex items-baseline space-x-4">
                                <NavLink to="/" exact>
                                    Home
                                </NavLink>
                                {auth.isLoggedIn && (
                                    <>
                                        <NavLink to="/budgets">Budgets</NavLink>

                                        <NavLink to="/transactions">
                                            Transactions
                                        </NavLink>

                                        <NavLink to="/categories">
                                            Categories
                                        </NavLink>
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                    <div className="md:hidden">
                        <button
                            className="text-gray-400 hover:text-white focus:outline-none focus:text-white"
                            onClick={toggleMobileMenu}
                        >
                            <svg
                                className="h-6 w-6"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d={
                                        isMobileMenuOpen
                                            ? "M6 18L18 6M6 6l12 12"
                                            : "M4 6h16M4 12h16M4 18h16"
                                    }
                                />
                            </svg>
                        </button>
                    </div>
                    <div className="hidden md:block">
                        {!auth.isLoggedIn && (
                            <Link
                                to="/login"
                                className="bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium"
                            >
                                Login
                            </Link>
                        )}

                        {auth.isLoggedIn && (
                            <Link
                                to="/"
                                className="bg-gray-900 text-white px-3 py-2 rounded-md text-sm font-medium"
                                onClick={() => {
                                    auth.logout();
                                }}
                            >
                                Logout
                            </Link>
                        )}
                    </div>
                </div>
            </div>
            {isMobileMenuOpen && (
                <div className="md:hidden">
                    <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
                        <NavLink to="/" exact>
                            Home
                        </NavLink>
                        {auth.isLoggedIn && (
                            <>
                                <NavLink to="/budgets">Budgets</NavLink>

                                <NavLink to="/transactions">
                                    Transactions
                                </NavLink>

                                <NavLink to="/categories">Categories</NavLink>
                            </>
                        )}
                    </div>
                    <div className="pt-4 pb-3 border-t border-gray-700">
                        {!auth.isLoggedIn && (
                            <div className="flex items-center px-4">
                                <Link
                                    to="/login"
                                    className="text-gray-400 hover:text-white mr-4"
                                >
                                    <span className="text-sm font-medium">
                                        Login
                                    </span>
                                </Link>
                            </div>
                        )}
                        {auth.isLoggedIn && (
                            <div className="flex items-center px-4">
                                <Link
                                    to="/"
                                    className="text-gray-400 hover:text-white mr-4"
                                    onClick={() => {
                                        auth.logout();
                                    }}
                                >
                                    <span className="text-sm font-medium">
                                        Logout
                                    </span>
                                </Link>
                            </div>
                        )}
                    </div>
                </div>
            )}
        </nav>
    );
};
