import Centered from "./Centered";
import { NAVBAR_HEIGHT_PX } from "../constants";
import { Outlet } from "@solidjs/router";

export const metadata = {
  title: "Boards",
  description: "Collaborate with your team",
};

export default function AuthLayout() {
  return (
    <div style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }}>
      <Centered>
        <Outlet/>
      </Centered>
    </div>
  );
}
