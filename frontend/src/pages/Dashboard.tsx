import { A } from "@solidjs/router";
import { createSignal, type Component, createEffect } from "solid-js";
import { COOKIE_NAME_JWT_TOKEN } from "../constants";
import Cookies from "js-cookie";

const Price: Component = () => {
  const [price, setPrice] = createSignal<number | null>(null);

  createEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/price");
    // Set up event listeners
    ws.addEventListener("open", () => {
      console.log("WebSocket connection opened");
    });

    ws.addEventListener("message", (event) => {
      const receivedMessage = event.data;
      console.log("Received:", receivedMessage);

      // Update state with the latest message
      setPrice(receivedMessage);
    });

    // Clean up: Close the WebSocket connection when this effect is destroyed
    return () => {
      console.log("WebSocket connection closed");
      ws.close();
    };
  });
  return (
    <div class="flex bg-sky-700 justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      <header class="gap-y-3 text-sky-200 flex items-center flex-col">
        <button
          class="bg-sky-200 hover:bg-sky-300 text-sky-700 font-semibold py-2 px-4 border border-sky-500 rounded shadow"
          onClick={() => {
            Cookies.remove(COOKIE_NAME_JWT_TOKEN);
            window.location.href = "/";
          }}
        >
          Sign out
        </button>
        <A href="/" class="hover:underline">
          Home
        </A>
        <p class=" italic underline-offset-4">
          market price: {price() ?? "price data not available"}
        </p>
      </header>
    </div>
  );
};

export default Price;
