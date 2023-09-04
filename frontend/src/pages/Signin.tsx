import SignInForm from "../components/SigninForm";
import Redirect from "../components/Redirect";

export default function SignInPage() {
  return (
    <div class="card card-bordered shadow-xl w-96 bg-white">
      <div class="card-body space-y-3">
        <div class="card-title text-2xl">Sign in</div>
        <SignInForm />
        <div class="divider">OR</div>
        <Redirect
          url="/signup"
          redirectText="Don't have an account?"
          buttonText="Sign up"
        />
      </div>
    </div>
  );
}
