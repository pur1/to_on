import { LogIn, UserRoundPlus } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function Header() {
  const navigate = useNavigate();

  return (
    <div className="flex bg-gray-100">
      <div className="flex flex-1" />
      <div className="flex">
        <button
          className="flex flex-col items-center py-1 px-3 w-20 text-xs font-bold bg-gray-200 hover:bg-gray-300"
          onClick={() => {
            navigate("/register");
          }}
        >
          <UserRoundPlus className="m-1" />
          Register
        </button>
        <button
          className="flex flex-col items-center py-1 px-3 w-20 text-xs font-bold bg-gray-200 hover:bg-gray-300"
          onClick={() => {
            navigate("/login");
          }}
        >
          <LogIn className="m-1" />
          Login
        </button>
      </div>
    </div>
  );
}
