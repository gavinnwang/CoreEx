import { A } from "@solidjs/router";
import { LogoIcon } from "./icons/LogoIcon";

export default function Logo({ className }: { className?: string }) {
  return (
    <A href="/" class={className}>
      <div class="flex gap-x-1 items-center italic">
        <LogoIcon />
        CoreEx
      </div>
    </A>
  );
}
