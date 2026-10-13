import GoogleIcon from "../components/icons/GoogleIcon";
import FacebookIcon from "../components/icons/FacebookIcon";
import { useState } from "react";

export default function RegisterPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <section className="flex flex-col gap-4 p-6 bg-white border border-gray-300 rounded w-full max-w-sm">
        <div className="flex-1 flex flex-col gap-3">
          <p className="mb-6 text-gray-600 text-left">Register</p>
          <input
            placeholder="Email"
            className="border border-gray-800 py-2 px-4 focus:outline-none focus:ring-1 focus:ring-gray-900"
          />
          <input
            placeholder="Username"
            className="border border-gray-800 py-2 px-4 focus:outline-none focus:ring-1 focus:ring-gray-900"
          />
          <div className="relative">
            <input
              placeholder="Password"
              type={showPassword ? "text" : "password"}
              className="w-full border border-gray-800 py-2 pl-4 pr-9 focus:outline-none focus:ring-1 focus:ring-gray-900"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              aria-label={showPassword ? "Hide password" : "Show password"}
              className="absolute text-xs hover:underline inset-y-0 right-0 flex items-center px-4 text-gray-500 hover:text-gray-800"
            >
              {showPassword ? <span>hide</span> : <span>show</span>}
            </button>
          </div>
          <div className="relative">
            <input
              placeholder="Confirm Password"
              type={showConfirmPassword ? "text" : "password"}
              className="w-full border border-gray-800 py-2 pl-4 pr-9 focus:outline-none focus:ring-1 focus:ring-gray-900"
            />
            <button
              type="button"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              aria-label={
                showConfirmPassword ? "Hide password" : "Show password"
              }
              className="absolute text-xs hover:underline inset-y-0 right-0 flex items-center px-4 text-gray-500 hover:text-gray-800"
            >
              {showConfirmPassword ? <span>hide</span> : <span>show</span>}
            </button>
          </div>
          <button className="px-6 py-3 font-bold text-white bg-orange-400 rounded hover:bg-orange-300 border border-orange-300">
            Register
          </button>
        </div>

        <div className="flex items-center gap-2 text-sm text-gray-400">
          <div className="h-px flex-1 bg-gray-300" />
          <span>OR</span>
          <div className="h-px flex-1 bg-gray-300" />
        </div>

        <div className="flex-1 flex flex-col gap-3 items-center justify-center">
          <button className="flex w-full h-12 items-center text-left gap-2 px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">
            <GoogleIcon className="w-6 h-6" />
            Login via Google
          </button>
          <button className="flex w-full h-12 items-center text-left gap-2 px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">
            <FacebookIcon className="w-6 h-6" />
            Login via Facebook
          </button>
        </div>
      </section>

      <span className="mt-2 text-sm text-gray-500">
        Already have an account?{" "}
        <a href="/login" className="text-sm text-blue-500 hover:underline">
          Login
        </a>
      </span>

      <a href="/" className="mt-2 text-sm text-blue-500 hover:underline">
        Return to Home
      </a>
    </div>
  );
}
