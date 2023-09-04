import { A } from "@solidjs/router";
import Cookies from "js-cookie";
import { COOKIE_NAME_JWT_TOKEN } from "./constants";
import { Toaster } from "solid-toast";

function getToken(): string | undefined {
  const jwtToken = Cookies.get(COOKIE_NAME_JWT_TOKEN);
  return jwtToken;
}
export default () => {
  const token = getToken();
  return (
    <div class="flex  justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      <Toaster />
      <h1 class="text-4xl font-semibold ">Exchange</h1>
      <A href={ token ? "/dashboard" : "signup" }class="hover:underline">
        Enter the exchange
      </A>
    </div>
  );
};
