import { A, Outlet } from "@solidjs/router";
import { createResource } from "solid-js";
import { User, getUserByJwt } from "../api/user";
import { COOKIE_NAME_JWT_TOKEN } from "../constants";
import Cookies from "js-cookie";
import Navbar from "./NavBar";
import { Toaster } from "solid-toast";

function getToken(): string | undefined {
  const jwtToken = Cookies.get(COOKIE_NAME_JWT_TOKEN);
  return jwtToken;
}

async function getUser(token: string | undefined): Promise<User | null> {
  if (token) {
    try {
      const user = await getUserByJwt(token);
      console.log("user", user);
      return user;
    } catch (e) {
      return null;
    }
  }
  return null;
}

export default () => {
  const token = getToken();
  const [user] = createResource(token, getUser);
  return (
    <div>
      <Toaster />
      <Navbar user={user} />
      <Outlet />
    </div>
  );
};
