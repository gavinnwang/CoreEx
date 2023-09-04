import { A, useNavigate } from "@solidjs/router";
import { createSignal, type Component, createEffect } from "solid-js";
import { COOKIE_NAME_JWT_TOKEN, NAVBAR_HEIGHT_PX } from "../constants";
import Cookies from "js-cookie";

const Price: Component = () => {
  const [price, setPrice] = createSignal<number | null>(null);

  const navigator = useNavigate();
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
    <div
      class="hero bg-base-200"
      style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }}
    >
      <header class="">
        <p class=" italic underline-offset-4">
          market price: {price() ?? "price data not available"}
        </p>
      </header>
    </div>
  );
};

export default Price;
