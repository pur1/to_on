import { Navigate, Outlet } from "react-router-dom";
import { useCurrentUser } from "../../lib/auth";

export default function GuestOnly() {
  const { data: user, isPending } = useCurrentUser();

  if (isPending) return null;
  if (user) return <Navigate to="/" replace />;
  return <Outlet />;
}
