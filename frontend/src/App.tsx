import { A } from "@solidjs/router";

export default () => {
  return (
    <div class="flex  justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      hi
      <A href="/signup">sign up / log in</A> 
      <A href="/price">price</A>
    </div>
  );
};
