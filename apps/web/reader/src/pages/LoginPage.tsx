import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useLogin } from "../lib/auth";
import AuthCard from "../components/auth/AuthCard";
import OrDivider from "../components/auth/OrDivider";
import SocialLoginButtons from "../components/auth/SocialLoginButtons";
import { PasswordInput, TextInput } from "@to-on/ui";

export default function LoginPage() {
  const navigate = useNavigate();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");

  const login = useLogin();

  return (
    <AuthCard
      className="max-w-3xl"
      footer={
        <>
          Don't have an account?{" "}
          <a href="/register" className="text-accent hover:underline">
            Register
          </a>
        </>
      }
    >
      <div className="flex-1 flex flex-col gap-3">
        <p className="mb-6 text-muted text-left">Login</p>
        <TextInput
          placeholder="Email or Username"
          value={identifier}
          onChange={(e) => setIdentifier(e.target.value)}
        />
        <PasswordInput
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {login.error && (
          <p className="text-left text-xs text-danger">{login.error.message}</p>
        )}
        <button
          type="button"
          onClick={() =>
            login.mutate(
              { identifier, password },
              { onSuccess: () => navigate("/") },
            )
          }
          disabled={login.isPending}
          className="px-6 py-3 font-bold text-primary-foreground bg-primary rounded hover:bg-primary-hover border border-primary disabled:opacity-50"
        >
          Login
        </button>
      </div>

      <OrDivider vertical />

      <SocialLoginButtons />
    </AuthCard>
  );
}
