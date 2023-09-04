import { JSXElement } from "solid-js";

export default function Centered({ children }: { children: JSXElement }) {
    return <div class="flex items-center justify-center h-full">{children}</div>;
  }
  