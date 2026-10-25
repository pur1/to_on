import type { ReactNode } from "react";

type AuthCardProps = {
  children: ReactNode;
  footer: ReactNode;
  className?: string;
};

export default function AuthCard({
  children,
  footer,
  className = "",
}: AuthCardProps) {
  return (
    <div
      data-palette="classic"
      className="flex flex-col items-center justify-center min-h-screen bg-background text-foreground"
    >
      <section
        className={`flex gap-4 p-6 bg-surface border border-border rounded w-full ${className}`}
      >
        {children}
      </section>

      <span className="mt-2 text-sm text-muted">{footer}</span>

      <a href="/" className="mt-2 text-sm text-accent hover:underline">
        Return to Home
      </a>
    </div>
  );
}
