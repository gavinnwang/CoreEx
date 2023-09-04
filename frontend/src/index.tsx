import { render } from "solid-js/web";
import App from "./App";
import "./assets/global.css";
import { Route, Router, Routes } from "@solidjs/router";
import { lazy } from "solid-js";
const root = document.getElementById("root");

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    "Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?"
  );
}

const Signup = lazy(() => import("./pages/Signup"));
const Dashboard = lazy(() => import("./pages/Dashboard"));
const Login = lazy(() => import("./pages/Login"));

render(
  () => (
    <Router>
      <Routes>
        <Route path="/" component={App} />
        <Route path="/signup" component={Signup} />
        <Route path="/login" component={Login} />
        <Route path="/dashboard" component={Dashboard} />
      </Routes>
    </Router>
  ),
  root!
);
