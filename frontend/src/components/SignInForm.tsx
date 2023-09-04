import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";
import { setSessionCookie } from "../lib/cookie";
import toast from "solid-toast";
import { signin } from "../api/auth";
import { APIError } from "../api";

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
      toast.success("Logged in successfully");
      navigate("/");
      window.location.reload();
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
