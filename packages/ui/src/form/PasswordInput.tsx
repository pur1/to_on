import { useState, type ComponentProps } from "react";
import TextInput from "./TextInput";

type PasswordInputProps = Omit<ComponentProps<typeof TextInput>, "type">;

export default function PasswordInput(props: PasswordInputProps) {
  const [show, setShow] = useState(false);

  return (
    <TextInput {...props} type={show ? "text" : "password"} className="pr-9">
      <button
        type="button"
        onClick={() => setShow(!show)}
        aria-label={show ? "Hide password" : "Show password"}
        className="absolute text-xs hover:underline inset-y-0 right-0 flex items-center px-4 text-muted hover:text-foreground"
      >
        {show ? "hide" : "show"}
      </button>
    </TextInput>
  );
}
