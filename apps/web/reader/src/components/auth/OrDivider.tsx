type OrDividerProps = {
  vertical?: boolean;
};

export default function OrDivider({ vertical = false }: OrDividerProps) {
  const line = vertical ? "w-px flex-1 bg-border" : "h-px flex-1 bg-border";

  return (
    <div
      className={`flex items-center gap-2 text-sm text-muted ${vertical ? "flex-col" : ""}`}
    >
      <div className={line} />
      <span>OR</span>
      <div className={line} />
    </div>
  );
}
