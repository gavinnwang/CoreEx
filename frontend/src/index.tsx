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

const Auth = lazy(() => import("./pages/Auth"));
const Price = lazy(() => import("./pages/Price"));

render(
  () => (
    <Router>
      <Routes>
        <Route path="/" component={App} />
        <Route path="/login" component={Auth} />
        <Route path="/price" component={Price} />
      </Routes>
    </Router>
  ),
  root!
);
