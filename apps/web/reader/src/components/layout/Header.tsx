import { useNavigate } from "react-router-dom";
import { LanguageSelect, useTranslation } from "@to-on/i18n";
import { useCurrentUser, useLogout } from "../../lib/auth";
import AuthButton from "./AuthButton";
import NavItem from "./NavItem";
import UserIcon from "../icons/UserIcon";

export default function Header() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { data: user } = useCurrentUser();
  const logout = useLogout();

  return (
    <header className="relative flex bg-header border-b border-border px-4 py-2">
      <div className="relative flex-1 mx-auto max-w-6xl">
        <div className="flex items-center justify-end gap-3 text-xs text-header-foreground">
          <NavItem
            className="mr-auto"
            label="Creator Studio"
            href={import.meta.env.VITE_STUDIO_URL}
            newTab
          />
          <LanguageSelect className="bg-header" />

          <div className="flex items-center">
            {user ? (
              <>
                <AuthButton
                  label={t("logout")}
                  onClick={() => logout.mutate()}
                  className="hover:text-header-foreground-hover"
                />
              </>
            ) : (
              <>
                <AuthButton
                  label={t("register")}
                  onClick={() => navigate("/register")}
                  className="hover:text-header-foreground-hover"
                />
                <span aria-hidden className="mx-2 h-3 w-px bg-border" />
                <AuthButton
                  label={t("login")}
                  onClick={() => navigate("/login")}
                  className="hover:text-header-foreground-hover"
                />
              </>
            )}
          </div>
        </div>
        <div className="flex">
          {user ? (
            <span className="flex flex-col font-bold mt-3 ml-auto bg-background px-4 py-1 items-center">
              <span>{user.username}</span>
              <span className="text-xs">{user.role.toUpperCase()}</span>
            </span>
          ) : (
            <span className="flex mt-3 ml-auto font-bold bg-background px-4 py-1 items-center">
              <UserIcon className="size-6" />
              Guest
            </span>
          )}
        </div>
      </div>
    </header>
  );
}
