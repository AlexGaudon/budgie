import { Outlet } from "react-router-dom";

import { NavBar } from "./components/NavBar";

export const Layout = () => {
    return (
        <>
            <NavBar></NavBar>
            <div className="h-screen flex flec-col">
                <div
                    id="content"
                    className="flex-grow mx-auto p-8 w-full shadow rounded"
                >
                    <Outlet></Outlet>
                </div>
            </div>
        </>
    );
};
