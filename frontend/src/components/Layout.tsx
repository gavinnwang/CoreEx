import { Outlet } from "@solidjs/router";
import Navbar from "./NavBar";
import { Toaster } from "solid-toast";

export default () => {
  return (
    <div class="pt-16">
      <Navbar />
      <Toaster />
      <Outlet />
    </div>
  );
};
