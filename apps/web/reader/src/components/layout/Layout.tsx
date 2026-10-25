import { Outlet } from "react-router-dom";
import Header from "./Header";

export default function Layout() {
  return (
    <>
      <Header />
      <main className="flex-1 w-full max-w-6xl mx-auto px-4">
        <Outlet />
      </main>
    </>
  );
}
