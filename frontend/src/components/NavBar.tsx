import { A } from "@solidjs/router";
import { User } from "../api/user";
import { NAVBAR_HEIGHT_PX } from "../constants";
import WidthContainer from "./WidthContainer";
import { Resource } from "solid-js";
import Avatar from "./Avatar";
import AccountMenu from "./AccountMenu";
import LogoText from "./LogoText";
import { Logo } from "./icons/Logo";

export default function Navbar({ user }: { user: Resource<User | null> }) {
  return (
    <div
      class="navbar fixed top-0 left-0 w-full bg-white border-b "
      style={{ height: NAVBAR_HEIGHT_PX, "z-index": 10002 }}
    >
      <WidthContainer>
        <div class="flex justify-between items-center w-full">
          <div class="flex gap-x-1 items-center italic">
            <Logo />
            <LogoText className="font-semibold text-xl" />
          </div>
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
