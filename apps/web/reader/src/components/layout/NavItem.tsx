type NavItemProps = {
  label: string;
  className?: string;
} & (
  | { href: string; newTab?: boolean; onClick?: never }
  | { onClick: () => void; href?: never; newTab?: never }
);

export default function NavItem({
  label,
  href,
  newTab = false,
  onClick,
  className = "",
}: NavItemProps) {
  const classes = `flex flex-col my-1 ${className}`;

  if (href) {
    return (
      <a
        href={href}
        className={classes}
        {...(newTab && { target: "_blank", rel: "noopener noreferrer" })}
      >
        {label}
      </a>
    );
  }

  return (
    <button className={classes} onClick={onClick}>
      {label}
    </button>
  );
}
