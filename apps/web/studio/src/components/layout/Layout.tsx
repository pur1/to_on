import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";

export default function Layout() {
  return (
    <div className="flex flex-1">
      <Sidebar />
      <main className="flex-1 p-6 text-left">
        <Outlet />
      </main>
    </div>
  );
}
