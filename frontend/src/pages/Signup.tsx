import Redirect from "../components/Redirect";
import SignUpForm from "../components/SignUpForm";

export default function SignUpPage() {
  return (
    <div class="card card-bordered shadow-xl w-96 bg-white">
      <div class="card-body space-y-3">
        <div class="card-title text-2xl">Sign up</div>
        <SignUpForm />
        <div class="divider">OR</div>
        <Redirect
          url="/signin"
          redirectText="Already have an account?"
          buttonText="Sign in"
        />
      </div>
    </div>
  );
}
