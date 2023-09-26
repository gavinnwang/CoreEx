import {
  createSignal,
  type Component,
  createEffect,
  createResource,
} from "solid-js";
import { BASE_URL, NAVBAR_HEIGHT_PX } from "../constants";
import PlaceOrderForm from "../components/PlaceOrderForm";
import SymbolInfoTable from "../components/SymbolInfoTable";
import toast from "solid-toast";
import { getPriceData } from "../api/marketData";
import { SolidApexCharts } from "solid-apexcharts";



const Price: Component = () => {
  const [symbolInfo, setSymbolInfo] = createSignal<SymbolInfo | null>(null);
  // const [candleData, setCandleData] = createSignal<CandleDataUpdate | null>(
  //   null
  // );
  const [fetchErrorMsg, setFetchErrorMsg] = createSignal<string | null>(null);

  const [graphData, setGraphData] = createSignal<ApexGraphData | null>(null);

  const symbol = "AAPL";

  const [priceData] = createResource(symbol, getPriceHistoryData);

  const convertDataToGraphFormat = (
  data: CandleDataPoint[]
): { data: GraphPriceDataPoint[] } [] => {
  return [
    {
      data: data.map((item) => ({
        x: new Date(item.recorded_at),
        y: [item.open, item.high, item.low, item.close] as [
          number,
          number,
          number,
          number
        ],
      })),
    },
  ];
};

async function getPriceHistoryData(symbol: string): Promise<ApexGraphData > {
  try {
    const data = await getPriceData(symbol);
    const graphData = convertDataToGraphFormat(data);
    setGraphData(graphData);
    return graphData;
  } catch (e) {
    console.error(e);
    throw e;
  }
}


  // createEffect(() => {


  //   const gd = graphData();
  //   const cd = candleData();
  //   if (cd && gd) {
  //     if (cd.new_candle) {
  //        gd[0].data.push({
  //         x: new Date(cd.recorded_at),
  //         y: [
  //           cd.open,
  //           cd.high,
  //           cd.low,
  //           cd.close,
  //         ] ,
  //       });
        
  //     } else {
  //       gd[0].data[gd[0].data.length - 1] = {
  //         x: new Date(cd.recorded_at),
  //         y: [
  //           cd.open,
  //           cd.high,
  //           cd.low,
  //           cd.close,
  //         ] ,
  //       };
  //     }
  //     toast.success(gd[0].data.length)
  //     setGraphData(gd);
  //   }
  // });

  const [needNew, setNeedNew] = createSignal<boolean>(true);

  createEffect(() => {
    const url = `ws://${BASE_URL}/ws`;
    const ws = new WebSocket(url);

    ws.addEventListener("open", () => {
      console.log("WebSocket connection opened");
      const payload: ParamsStreamPrice = {
        event: "exchange.stream_info",
        params: {
          symbol: symbol,
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
      if (resData.success && resData.result) {
        setSymbolInfo(resData.result);
        // setCandleData(resData.result.candle_data);
        const cd = resData.result.candle_data
        if (resData.result.candle_data) {
          const gd = [...graphData() as ApexGraphData]

          if (needNew() && !cd.new_candle) {
            gd[0].data.push({
              x: new Date(cd.recorded_at),
              y: [
                cd.open,
                cd.high,
                cd.low,
                cd.close,
              ] ,
            });
            setNeedNew(false);
          }

          if (needNew() && cd.new_candle) {
            gd[0].data.push({
              x: new Date(cd.recorded_at),
              y: [
                cd.open,
                cd.high,
                cd.low,
                cd.close,
              ] ,
            });
            setNeedNew(true);
          }

          if (!needNew() && !cd.new_candle) {
            gd[0].data[gd[0].data.length - 1] = {
              x: new Date(cd.recorded_at),
              y: [
                cd.open,
                cd.high,
                cd.low,
                cd.close,
              ] ,
            };

          } 

          if (!needNew() && cd.new_candle) {
            gd[0].data[gd[0].data.length - 1] = {
              x: new Date(cd.recorded_at),
              y: [
                cd.open,
                cd.high,
                cd.low,
                cd.close,
              ] ,
            };
            setNeedNew(true);
          }

          // if (cd.new_candle && !createNewCandle()) {
          //    gd[0].data.push({
          //     x: new Date(cd.recorded_at),
          //     y: [
          //       cd.open,
          //       cd.high,
          //       cd.low,
          //       cd.close,
          //     ] ,
          //   });

            
          // } else if (cd.new_candle && createNewCandle()) {
          
          // else if (!cd.new_candle && createNewCandle()) {
          //   gd[0].data[gd[0].data.length - 1] = {
          //     x: new Date(cd.recorded_at),
          //     y: [
          //       cd.open,
          //       cd.high,
          //       cd.low,
          //       cd.close,
          //     ] ,
          //   };
          // }
          setGraphData(gd);
        }
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
          {priceData.loading ? (
            <div>Loading...</div>
          ) : priceData.error ? (
            <div>Error: {priceData.error.message}</div>
          ) : graphData() &&  (
            (
              <SolidApexCharts
                width="800"
                type="candlestick"
                options={{
                  chart: {
                    id: "solidchart-example",
                  },
                }}
                series={graphData() as ApexGraphData}
              />
            )
          )}
        </div>
      </div>
    </div>
  );
};

export default Price;
