import { A } from "@solidjs/router";
import { token } from "./store";
import { NAVBAR_HEIGHT_PX } from "./constants";


export default () => {
  return (
    <div
      class="hero bg-base-200"
      style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }}
    >
      <A
        href={token() ? "/dashboard" : "/signup"}
        class="btn btn-secondary btn-outline"
      >
        Enter the exchange
      </A>
    </div>
  );
};
