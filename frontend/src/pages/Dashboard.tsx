import { createSignal, type Component, createEffect } from "solid-js";
import { BASE_URL, NAVBAR_HEIGHT_PX } from "../constants";
import PlaceOrderForm from "../components/PlaceOrderForm";
import CandleGraph from "../components/CandleGraph";

const Price: Component = () => {
  const [price, setPrice] = createSignal<number | null>(null);
  const [fetchErrorMsg, setFetchErrorMsg] = createSignal<string | null>(null);

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
      setFetchErrorMsg("Something went wrong.");
    });

    ws.addEventListener("message", (event) => {
      const res = event.data;
      const resData: ResponseGetMarketPrice = JSON.parse(res);
      if (resData.success) {
        setPrice(resData.result ? resData.result.price : null);
      } else {
        setFetchErrorMsg(resData.error_message ?? "Something went wrong.");
        setPrice(null);
      }
    });

    return () => {
      console.log("WebSocket connection closed");
      ws.close();
    };
  });

  return (
    <div style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }} class="p-5">
      <div class="text-xl font-semibold py-4">
        {fetchErrorMsg()
          ? fetchErrorMsg()
          : price() !== null
          ? `AAPL price: \$${price()}`
          : "loading..."}
      </div>
      <PlaceOrderForm />
      <CandleGraph />
    </div>
  );
};

export default Price;
