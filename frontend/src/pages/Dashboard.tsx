import { createSignal, type Component, createEffect } from "solid-js";
import { BASE_URL, NAVBAR_HEIGHT_PX } from "../constants";
import PlaceOrderForm from "../components/PlaceOrderForm";
import CandleGraph from "../components/CandleGraph";
import SymbolInfoTable from "../components/SymbolInfoTable";
import toast from "solid-toast";

const Price: Component = () => {
  const [symbolInfo, setSymbolInfo] = createSignal<SymbolInfo | null>(null);
  const [fetchErrorMsg, setFetchErrorMsg] = createSignal<string | null>(null);

  createEffect(() => {
    const url = `ws://${BASE_URL}/ws`;
    const ws = new WebSocket(url);

    ws.addEventListener("open", () => {
      console.log("WebSocket connection opened");
      const payload: ParamsStreamPrice = {
        event: "exchange.stream_info",
        params: {
          symbol: "AAPL",
        },
      };
      ws.send(JSON.stringify(payload));
    });
    ws.addEventListener("error", (error) => {
      console.error("WebSocket Error:", error);
      setFetchErrorMsg("Something went wrong.");
      toast.error("Something went wrong.");
    });

    ws.addEventListener("message", (event) => {
      const res = event.data;
      const resData: WSResponseGetSymbolInfo = JSON.parse(res);
      // console.log(resData)
      if (resData.success && resData.result) {
        setSymbolInfo(resData.result);
      } else {
        const errMsg = resData.error_message ?? "Something went wrong.";
        setFetchErrorMsg(errMsg);
        toast.error(errMsg);
      }
    });

    return () => {
      console.log("WebSocket connection closed");
      ws.close();
    };
  });

  return (
    <div
      style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }}
      class="py-5 px-10"
    >
      <div class="flex flex-col gap-y-5 ">
        <PlaceOrderForm />
        <div class="flex flex-col gap-y-5 md:flex-row md:gap-x-5">
          <SymbolInfoTable
            symbolInfo={symbolInfo}
            fetchErrorMsg={fetchErrorMsg}
          />
          <CandleGraph />
        </div>
      </div>
    </div>
  );
};

export default Price;
