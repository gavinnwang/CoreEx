import { A, useNavigate } from "@solidjs/router";
import { createSignal, type Component } from "solid-js";
import { COOKIE_NAME_JWT_TOKEN } from "../constants";
import Cookies from "js-cookie";
import { APIError } from "../api";
import toast, { Toaster } from "solid-toast";
import { login } from "../api/auth";

const Login: Component = () => {
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");

  const navigate = useNavigate();
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const { token } = await login({
        email: email(),
        password: password(),
      });
      const expirationDate = new Date();
      expirationDate.setDate(expirationDate.getDate() + 30);

      toast.success("Logged in successfully");
      Cookies.set(COOKIE_NAME_JWT_TOKEN, token, {
        path: "/",
        expires: expirationDate,
        httpOnly: true,
        secure: true,
      });

      // Refresh and navigate (assuming you have some refresh mechanism)
      // router.refresh();
      navigate("/price");
    } catch (error) {
      if (error instanceof APIError) {
        toast.error(error.message);
      } else {
        toast.error("An unknown error occurred");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div class="flex justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      <A class="hover:underline" href="/">
        home
      </A>
      Sign Up
      <Toaster />
      <form
        onSubmit={handleSubmit}
        class="flex flex-col gap-y-3 border rounded-md p-4 shadow-sm"
      >
        <input
          type="email"
          placeholder="Email"
          onInput={(e) => setEmail(e.currentTarget.value)}
          class="border-b focus:outline-none"
          required
        />
        <input
          type="password"
          placeholder="Password"
          onInput={(e) => setPassword(e.currentTarget.value)}
          class="border-b focus:outline-none"
          required
        />
        <button class="hover:underline" type="submit" disabled={isLoading()}>
          {isLoading() ? "Loading..." : "Log In"}
        </button>
      </form>
      <A href="/signup" class="hover:underline">
        Don&apos;t have an account? Sign up!{" "}
      </A>
    </div>
  );
};

export default Login;
