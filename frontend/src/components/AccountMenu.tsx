import { A, useNavigate } from "@solidjs/router";
import Cookies from "js-cookie";
import { COOKIE_NAME_JWT_TOKEN } from "../constants";
import { User } from "../api/user";
import { JSX } from "solid-js";
import { ChevDown } from "./icons/ChevDown";

export default function AccountMenu({
  user,
  avatar,
}: {
  user: User;
  avatar: JSX.Element;
}) {
  const navigator = useNavigate();

  const handleLogout = () => {
    Cookies.remove(COOKIE_NAME_JWT_TOKEN);
    navigator("/");
    window.location.reload();
  };

  return (
    <div class="dropdown dropdown-end">
      <div tabIndex={0} class="btn btn-ghost normal-case ">
        <div class="w-10">{avatar}</div>
        <span>{user?.name}</span>
        <ChevDown />
      </div>
      <div class="right-0 mt-3 p-1 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-36">
        <ul class="menu menu-compact gap-1 p-3">
          <li>
            <A
              href="/dashboard"
              replace={true}
              class="flex items-center justify-between"
            >
              Dashboard
            </A>
          </li>
          <li>
            <button onClick={handleLogout}>Logout</button>
          </li>
        </ul>
      </div>
    </div>
  );
}
