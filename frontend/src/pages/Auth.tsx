import { useNavigate } from "@solidjs/router";
import { createSignal, type Component } from "solid-js";
import { createUser } from "../api/user";
import { COOKIE_NAME_JWT_TOKEN } from "../constants";
import Cookies from "js-cookie";
import { APIError } from "../api";
import toast, { Toaster } from "solid-toast";

const App: Component = () => {
  const [name, setName] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");

  const navigate = useNavigate();
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const { jwt_token: token } = await createUser({
        email: email(),
        password: password(),
        name: name(),
      });
      toast.success("Successfully created user " + name() + "!");
      const expirationDate = new Date();
      expirationDate.setDate(expirationDate.getDate() + 30);

      Cookies.set(COOKIE_NAME_JWT_TOKEN, token, {
        path: "/",
        expires: expirationDate,
      });

      // Refresh and navigate (assuming you have some refresh mechanism)
      // router.refresh();
      // navigate("/");
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
      Sign Up
      <Toaster />
      <form
        onSubmit={handleSubmit}
        class="flex flex-col gap-y-3 border rounded-md p-4 shadow-sm"
      >
        <input
          type="text"
          placeholder="Name"
          onInput={(e) => setName(e.currentTarget.value)}
          class="border-b focus:outline-none"
        />
        <input
          type="email"
          placeholder="Email"
          onInput={(e) => setEmail(e.currentTarget.value)}
          class="border-b focus:outline-none"
        />
        <input
          type="password"
          placeholder="Password"
          onInput={(e) => setPassword(e.currentTarget.value)}
          class="border-b focus:outline-none"
        />
        <button type="submit" disabled={isLoading()}>
          {isLoading() ? "Loading..." : "Sign Up"}
        </button>
      </form>
    </div>
  );
};

export default App;
