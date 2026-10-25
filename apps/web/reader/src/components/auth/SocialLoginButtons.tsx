import GoogleIcon from "../icons/GoogleIcon";
import FacebookIcon from "../icons/FacebookIcon";

const PROVIDERS = [
  { name: "Google", Icon: GoogleIcon },
  { name: "Facebook", Icon: FacebookIcon },
];

export default function SocialLoginButtons() {
  return (
    <div className="flex-1 flex flex-col gap-3 items-center justify-center">
      {PROVIDERS.map(({ name, Icon }) => (
        <button
          key={name}
          className="flex w-full h-12 items-center text-left gap-2 px-4 py-2 border border-border rounded hover:bg-surface-hover"
        >
          <Icon className="w-6 h-6" />
          Login via {name}
        </button>
      ))}
    </div>
  );
}
