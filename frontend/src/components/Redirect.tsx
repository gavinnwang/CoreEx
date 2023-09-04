import { A } from "@solidjs/router";

export default function Redirect({
  url,
  redirectText,
  buttonText,
}: {
  url: string;
  redirectText: string;
  buttonText: string;
}) {
  return (
    <div class="space-x-2">
      <span class="text-xs">{redirectText}</span>
      <A href={url} class="font-bold text-xs">
        {buttonText}
      </A>
    </div>
  );
}
