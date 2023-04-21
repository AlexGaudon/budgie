import { Link, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

interface NavLinkProps {
    to: string;
    exact?: boolean;
    children: React.ReactNode;
}

const NavLink: React.FC<NavLinkProps> = ({ to, exact = false, children }) => {
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

export const NavBar: React.FC = () => {
    let auth = useAuth();
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
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                    <div>
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
        </nav>
    );
};
