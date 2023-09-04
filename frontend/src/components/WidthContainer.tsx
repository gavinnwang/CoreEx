import { JSXElement } from "solid-js";

export default function WidthContainer (props: { children: JSXElement; class?: string })  {
  return (
    <div class={`max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full ${props.class}`}>
      {props.children}
    </div>
  );
};
