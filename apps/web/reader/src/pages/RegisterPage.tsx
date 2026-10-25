import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRegister } from "../lib/auth";
import AuthCard from "../components/auth/AuthCard";
import OrDivider from "../components/auth/OrDivider";
import SocialLoginButtons from "../components/auth/SocialLoginButtons";
import { PasswordInput, TextInput } from "@to-on/ui";

export default function RegisterPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const [emailError, setEmailError] = useState("");
  const [confirmError, setConfirmError] = useState("");

  const register = useRegister();

  const handleRegister = () => {
    const nextEmailError = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
      ? ""
      : "Invalid email format";
    const nextConfirmError =
      password !== confirmPassword ? "Passwords do not match" : "";

    setEmailError(nextEmailError);
    setConfirmError(nextConfirmError);

    if (nextEmailError || nextConfirmError) return;

    register.mutate(
      { email, username, password },
      { onSuccess: () => navigate("/") },
    );
  };

  return (
    <AuthCard
      className="flex-col max-w-sm"
      footer={
        <>
          Already have an account?{" "}
          <a href="/login" className="text-accent hover:underline">
            Login
          </a>
        </>
      }
    >
      <div className="flex-1 flex flex-col gap-4">
        <p className="mb-6 text-muted text-left">Register</p>
        <TextInput
          placeholder="Email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={emailError}
        />
        <TextInput
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <PasswordInput
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <PasswordInput
          placeholder="Confirm Password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          error={confirmError}
        />
        {register.error && (
          <p className="text-left text-xs text-danger">
            {register.error.message}
          </p>
        )}
        <button
          type="button"
          onClick={handleRegister}
          disabled={register.isPending}
          className="px-6 py-3 font-bold text-secondary-foreground bg-secondary rounded hover:bg-secondary-hover border border-secondary disabled:opacity-50"
        >
          Register
        </button>
      </div>

      <OrDivider />

      <SocialLoginButtons />
    </AuthCard>
  );
}
