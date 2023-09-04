import { A, useNavigate } from "@solidjs/router";
import { createSignal, type Component } from "solid-js";
import { createUser } from "../api/user";
import { APIError } from "../api";
import toast from "solid-toast";
import { setSessionCookie } from "../lib/cookie";
import { setToken } from "../store";

export default function SignUpForm() {
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

      setSessionCookie(token);
      setToken(token);
      toast.success("Successfully created user " + name() + "!");

      navigate("/");
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
    <form onSubmit={handleSubmit} class="space-y-4">
      <div class="form-control">
        <label class="label">
          <span class="label-text text-stone-600">Email</span>
        </label>
        <input
          type="email"
          id="email"
          autocomplete="email"
          class="input input-bordered w-full max-w-xs"
          onInput={(e) => setEmail(e.currentTarget.value)}
          required
        />
      </div>
      <div class="form-control">
        <label class="label">
          <span class="label-text text-stone-600">Name</span>
        </label>
        <input
          type="name"
          id="name"
          autocomplete="username"
          class="input input-bordered w-full max-w-xs"
          onInput={(e) => setName(e.currentTarget.value)}
          required
        />
      </div>
      <div class="form-control">
        <label class="label">
          <span class="label-text  text-stone-600">Password</span>
        </label>
        <input
          type="password"
          id="password"
          autocomplete="current-password"
          class="input input-bordered w-full max-w-x"
          onInput={(e) => setPassword(e.currentTarget.value)}
          required
        />
      </div>
      <div class="form-control mt-6">
        <button type="submit" class="btn btn-secondary" disabled={isLoading()}>
          {isLoading() ? "Signing in..." : "Sign up"}
        </button>
      </div>
    </form>
  );
}
