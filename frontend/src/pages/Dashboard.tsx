import { createSignal, type Component, createEffect } from "solid-js";
import { BASE_URL, NAVBAR_HEIGHT_PX } from "../constants";
import PlaceOrderForm from "../components/PlaceOrderForm";
import CandleGraph from "../components/CandleGraph";
import toast from "solid-toast";

const Price: Component = () => {
  const [price, setPrice] = createSignal<number | null>(null);
  const [fetchPriceError, setFetchPriceError] = createSignal(false);

  createEffect(() => {
    const url = `ws://${BASE_URL}/ws`;
    const ws = new WebSocket(url);

    ws.addEventListener("open", () => {
      console.log("WebSocket connection opened");
      const payload: ParamsStreamPrice = {
        event: "exchange.stream_price",
        params: {
          symbol: "AAPL",
        },
      };
      ws.send(JSON.stringify(payload));
    });

    ws.addEventListener("error", (error) => {
      console.error("WebSocket Error:", error);
      setFetchPriceError(true);
      toast.error("Error fetching price");
    });

    ws.addEventListener("message", (event) => {
      const res = event.data;
      const resData: ResponseGetMarketPrice = JSON.parse(res);
      setPrice(resData.result ? resData.result.price : null);
    });

    return () => {
      console.log("WebSocket connection closed");
      ws.close();
    };
  });
  return (
    <div
      class="bg-base-200"
      style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }}
    >
      <p class="italic">
        market price:
        {fetchPriceError() ? "something went wrong" : price() ?? "loading..."}
      </p>
      <PlaceOrderForm />
      <CandleGraph />
    </div>
  );
};

export default Price;
