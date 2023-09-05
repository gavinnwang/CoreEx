import { A } from "@solidjs/router";
import { User, getUserByJwt } from "../api/user";
import { COOKIE_NAME_JWT_TOKEN, NAVBAR_HEIGHT_PX } from "../constants";
import WidthContainer from "./WidthContainer";
import { createEffect } from "solid-js";
import Avatar from "./Avatar";
import AccountMenu from "./AccountMenu";
import LogoText from "./Logo";
import Cookies from "js-cookie";
import { setToken, setUser, token, user } from "../store";

function getToken(): string | undefined {
  const jwtToken = Cookies.get(COOKIE_NAME_JWT_TOKEN);
  return jwtToken;
}

async function getUser(token: string | undefined): Promise<User | null> {
  if (token) {
    try {
      const user = await getUserByJwt(token);
      return user;
    } catch (e) {
      return null;
    }
  }
  return null;
}

export default function Navbar() {
  createEffect(async () => {
    token()
    const t = getToken()
    if (!t) {
      setToken(undefined);
      setUser(null);
      return;
    }
    const currentUser = await getUser(t);
    if (currentUser) {
      setUser(currentUser);
      setToken(t)
    }
  });

  return (
    <div
      class="navbar fixed top-0 left-0 w-full bg-white shadow-md"
      style={{ height: NAVBAR_HEIGHT_PX }}
    >
      <WidthContainer>
        <div class="flex justify-between items-center w-full">
          <LogoText className="font-bold text-xl" />
          {user() ? (
            <div class="flex items-center">
              <AccountMenu user={user()!} avatar={<Avatar id={user()!.id} />} />
            </div>
          ) : (
            <AuthNav />
          )}
        </div>
      </WidthContainer>
    </div>
  );
}

function AuthNav() {
  return (
    <div class="space-x-2">
      <A href="/signin" class="btn btn-secondary btn-outline">
        Sign in
      </A>
      <A href="/signup" class="btn btn-primary">
        Sign up
      </A>
    </div>
  );
}
