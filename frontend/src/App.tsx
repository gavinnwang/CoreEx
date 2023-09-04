import { A } from "@solidjs/router";
import Cookies from "js-cookie";
import { COOKIE_NAME_JWT_TOKEN } from "./constants";
import { User, getUserByJwt } from "./api/user";
import { createResource } from "solid-js";

function getToken(): string | undefined {
  const jwtToken = Cookies.get(COOKIE_NAME_JWT_TOKEN);
  return jwtToken;
}

async function getUser(token: string | undefined): Promise<User | null> {
  if (token) {
    try {
      const user = await getUserByJwt(token);
      console.log("user", user)
      return user;
    } catch (e) {
      return null;
    }
  }
  return null;
}

export default () => {
  const token = getToken();

  const [user] = createResource(token, getUser)
  return (
    <div class="flex  justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      <h1 class="text-4xl font-semibold ">Exchange</h1>
      <A href={token ? "/dashboard" : "signup"} class="btn">
        Enter the exchange
      </A>
    </div>
  );
};
