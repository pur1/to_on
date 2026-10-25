import { useId, type InputHTMLAttributes, type ReactNode } from "react";

type TextInputProps = InputHTMLAttributes<HTMLInputElement> & {
  error?: string;
  children?: ReactNode;
};

export default function TextInput({
  error,
  children,
  className = "",
  ...props
}: TextInputProps) {
  const errorId = useId();

  return (
    <div className="relative">
      <input
        {...props}
        aria-invalid={!!error}
        aria-describedby={error ? errorId : undefined}
        className={`w-full border py-2 px-4 bg-input focus:outline-none focus:ring-1 ${
          error
            ? "border-danger focus:ring-danger"
            : "border-border focus:ring-accent"
        } ${className}`}
      />
      {children}
      {error && (
        <p
          id={errorId}
          className="absolute left-0 top-full text-left text-xs text-danger"
        >
          {error}
        </p>
      )}
    </div>
  );
}
