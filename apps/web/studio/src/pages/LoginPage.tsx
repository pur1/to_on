import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useLogin } from "../lib/auth";
import { LanguageSelect, useTranslation } from "@to-on/i18n";
import { PasswordInput, TextInput } from "@to-on/ui";

export default function LoginPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");

  const login = useLogin();

  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-background text-foreground">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          login.mutate(
            { identifier, password },
            { onSuccess: () => navigate("/dashboard", { replace: true }) },
          );
        }}
        className="flex flex-col gap-3 w-full max-w-sm p-6 bg-surface border border-border rounded"
      >
        <p className="mb-6 text-left font-bold">Creator Studio</p>
        <TextInput
          placeholder={t("identifier")}
          value={identifier}
          onChange={(e) => setIdentifier(e.target.value)}
        />
        <PasswordInput
          placeholder={t("password")}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {login.error && (
          <p className="text-left text-xs text-danger">{login.error.message}</p>
        )}
        <button
          type="submit"
          disabled={login.isPending}
          className="px-6 py-3 font-bold text-primary-foreground bg-primary rounded hover:bg-primary-hover border border-primary disabled:opacity-50"
        >
          {t("login")}
        </button>
      </form>

      <span className="mt-2 text-sm text-muted">
        {t("noAccount")}{" "}
        <a
          href={`${import.meta.env.VITE_READER_URL}/register`}
          className="text-accent hover:underline"
        >
          {t("register")}
        </a>
      </span>
      <LanguageSelect className="mt-4 text-sm text-muted bg-background" />
    </div>
  );
}
