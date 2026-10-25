import { NavLink } from "react-router-dom";
import { useTranslation } from "@to-on/i18n";

const NAV_ITEMS: {
  to: string;
  labelKey: string;
}[] = [
  { to: "/books", labelKey: "books" },
  { to: "/settings", labelKey: "settings" },
];

export default function Sidebar() {
  const { t } = useTranslation();

  return (
    <aside className="flex w-60 shrink-0 flex-col bg-header text-header-foreground border-r border-border">
      <div className="px-6 py-5 text-left text-lg font-bold">
        Creator Studio
      </div>
      <nav className="flex flex-col gap-1">
        {NAV_ITEMS.map(({ to, labelKey }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2 text-sm ${
                isActive
                  ? "bg-surface-hover text-foreground font-semibold"
                  : "hover:text-header-foreground-hover"
              }`
            }
          >
            {t(labelKey)}
          </NavLink>
        ))}
      </nav>
    </aside>
  );
}
