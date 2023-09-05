import { useNavigate } from "@solidjs/router";
import { createSignal, type Component, createEffect } from "solid-js";
import { BASE_URL, NAVBAR_HEIGHT_PX } from "../constants";

const Price: Component = () => {
  const [price, setPrice] = createSignal<number | null>(null);

  const navigator = useNavigate();
  createEffect(() => {
    const url = `ws://${BASE_URL}/ws`;
    const ws = new WebSocket(url);
    // Set up event listeners
    ws.addEventListener("open", () => {
      console.log("WebSocket connection opened");
      const payload : ParamsStreamPrice = {
        event: "exchange.stream_price",
        params: {
          symbol: "AAPL",
        },
      };
      ws.send(JSON.stringify(payload));
    });

    ws.addEventListener("message", (event) => {
      const res = event.data;
      const resData : ResponseGetMarketPrice = JSON.parse(res);

      // Update state with the latest message
      setPrice(resData.result ? resData.result.price : null);
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
