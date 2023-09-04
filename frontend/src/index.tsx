import { render } from "solid-js/web";
import App from "./App";
import "./assets/global.css";
import { Route, Router, Routes } from "@solidjs/router";
import { lazy } from "solid-js";
import Layout from "./components/Layout";
const root = document.getElementById("root");

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    "Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?"
  );
}

const SignUpPage = lazy(() => import("./pages/Signup"));
const Dashboard = lazy(() => import("./pages/Dashboard"));
const SignInPage = lazy(() => import("./pages/Signin"));
const AuthLayout = lazy(() => import("./components/AuthLayout"));

render(
  () => (
    <Router>
      <Routes>
        <Route path="/" component={Layout}>
          <Route path="/" component={App} />
          <Route path="/dashboard" component={Dashboard} />
          <Route path="/" component={AuthLayout}>
          <Route path="/signup" component={SignUpPage} />
          <Route path="/signin" component={SignInPage} />

          </Route>
        </Route>
      </Routes>
    </Router>
  ),
  root!
);
