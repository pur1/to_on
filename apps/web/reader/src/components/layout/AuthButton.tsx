type AuthButtonProps = {
  label: string;
  onClick: () => void;
  className?: string;
};

export default function AuthButton({
  label,
  onClick,
  className = "",
}: AuthButtonProps) {
  return (
    <button className={`flex flex-col my-1 ${className}`} onClick={onClick}>
      {label}
    </button>
  );
}
