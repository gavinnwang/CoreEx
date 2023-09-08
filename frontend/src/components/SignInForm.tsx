import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { setSessionCookie } from "../lib/cookie";
import { signin } from "../api/auth";
import { APIError } from "../api";
import { toast } from "solid-toast";
import { setToken } from "../store";

export default function SignInForm() {
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");

  const navigate = useNavigate();
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const { token } = await signin({ email: email(), password: password() });
      setSessionCookie(token);
      setToken(token)
      toast.success("Logged in successfully");
      navigate("/");
    } catch (error) {
      if (error instanceof APIError) {
        toast.error(error.message);
      } else {
        toast.error("Failed to connect to the server.");
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
        <button type="submit" class="btn btn-primary" disabled={isLoading()}>
          {isLoading() ? "Signing in..." : "Sign in"}
        </button>
      </div>
    </form>
  );
}
