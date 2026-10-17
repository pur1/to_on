import { LogIn, UserRoundPlus, Globe } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

export default function Header() {
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();

  const LANGS = [
    { code: "en", label: "English" },
    { code: "th", label: "ไทย" },
  ];

  return (
    <div className="flex flex-col bg-gray-100">
      <div className="bg-gray-100 p-1">
        <div className="flex items-center">
          <Globe size={15} />
          <select
            value={i18n.resolvedLanguage}
            onChange={(e) => i18n.changeLanguage(e.target.value)}
            className="text-xs focus:outline-none"
          >
            {LANGS.map(({ code, label }) => (
              <option key={code} value={code}>
                {label}
              </option>
            ))}
          </select>
        </div>
      </div>
      <div className="flex">
        <button
          className="flex flex-col items-center py-1 px-3 w-20 text-xs font-bold bg-gray-200 hover:bg-gray-300"
          onClick={() => {
            navigate("/register");
          }}
        >
          <UserRoundPlus className="m-1" />
          {t("register")}
        </button>
        <button
          className="flex flex-col items-center py-1 px-3 w-20 text-xs font-bold bg-gray-200 hover:bg-gray-300"
          onClick={() => {
            navigate("/login");
          }}
        >
          <LogIn className="m-1" />
          {t("login")}
        </button>
      </div>
    </div>
  );
}
