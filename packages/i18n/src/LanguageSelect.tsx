import { useTranslation } from "react-i18next";
import { ChevronDownIcon } from "@to-on/ui";

const LANGS = [
  { code: "en", label: "English" },
  { code: "th", label: "ภาษาไทย" },
];

// The native arrow sits at the far edge of the widest option, so it is
// replaced by an icon kept right next to the selected label.
export default function LanguageSelect({ className = "" }) {
  const { i18n } = useTranslation();

  return (
    <span className={`inline-grid items-center ${className}`}>
      <select
        value={i18n.resolvedLanguage}
        onChange={(e) => i18n.changeLanguage(e.target.value)}
        className="col-start-1 row-start-1 appearance-none field-sizing-content bg-transparent pr-4 leading-none cursor-pointer focus:outline-none"
      >
        {LANGS.map((lang) => (
          <option key={lang.code} value={lang.code}>
            {lang.label}
          </option>
        ))}
      </select>
      <ChevronDownIcon className="col-start-1 row-start-1 justify-self-end size-3 pointer-events-none" />
    </span>
  );
}
