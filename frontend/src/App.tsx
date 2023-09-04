import { A } from "@solidjs/router";

export default () => {
  return (
    <div class="flex  justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      hi
      <A href="/login">login</A> 
      <A href="/price">price</A>
    </div>
  );
};
