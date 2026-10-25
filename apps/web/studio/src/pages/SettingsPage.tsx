import { LanguageSelect, useTranslation } from "@to-on/i18n";

export default function SettingsPage() {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col gap-6 max-w-2xl">
      <h2>{t("settings")}</h2>

      <section className="flex items-center justify-between gap-6 p-6 bg-surface border border-border rounded">
        <div>
          <p className="font-semibold">{t("language")}</p>
          <p className="text-sm text-muted">{t("languageDescription")}</p>
        </div>
        <LanguageSelect className="px-3 py-2 text-sm bg-input border border-border rounded" />
      </section>
    </div>
  );
}
