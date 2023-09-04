import { A } from "@solidjs/router";

export default function LogoText({ className }: { className?: string }) {
  return (
    <A href="/" class={className}>
     CoreEx 
    </A>
  );
}
