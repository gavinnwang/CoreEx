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
import UserPrivateInfoPanel from "../components/UserPrivateInfoPanel";

const Price: Component = () => {
  const [symbolInfo, setSymbolInfo] = createSignal<SymbolInfo | null>(null);

  const [fetchErrorMsg, setFetchErrorMsg] = createSignal<string | null>(null);

  const [graphPriceData, setGraphPriceData] =
    createSignal<ApexGraphPriceData | null>(null);

  const [graphVolumeData, setGraphVolumeData] =
    createSignal<ApexGraphVolumeData | null>(null);

  const symbol = "AAPL";

  const [priceData] = createResource(symbol, getPriceHistoryData);

  const convertDataToPriceGraphFormat = (
    data: SymbolInfo[]
  ): { data: GraphPriceDataPoint[] }[] => {
    return [
      {
        data: data.map((cd) => ({
          x: cd.recorded_at,
          y: [cd.open, cd.high, cd.low, cd.close] as [
            number,
            number,
            number,
            number
          ],
        })),
      },
    ];
  };

  const convertDataToVolumeGraphFormat = (
    d: SymbolInfo[]
  ): ApexGraphVolumeData => {
    return [
      {
        name: "Ask Volume",
        data: d.map((cd) => ({
          x: cd.recorded_at,
          y: cd.ask_volume,
        })),
      },
      {
        name: "Bid Volume",
        data: d.map((cd) => ({
          x: cd.recorded_at,
          y: cd.bid_volume,
        })),
      },
    ];
  };

  async function getPriceHistoryData(symbol: string): Promise<SymbolInfo[]> {
    try {
      const data = await getPriceData(symbol);

      const graphPriceData = convertDataToPriceGraphFormat(data);
      const graphVolumeData = convertDataToVolumeGraphFormat(data);
      setGraphVolumeData(graphVolumeData);
      setGraphPriceData(graphPriceData);
      return data;
    } catch (e) {
      console.error(e);
      throw e;
    }
  }

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
      if (
        resData.success &&
        resData.result &&
        resData.event == "exchange.stream_info"
      ) {
        setSymbolInfo(resData.result);

        // setCandleData(resData.result.candle_data);
        const cd = resData.result;

        const gd = [...(graphPriceData() as ApexGraphPriceData)];
        const vd = [...(graphVolumeData() as ApexGraphVolumeData)];

        if (needNew() && !cd.new_candle) {
          // if (gd[0].data.length > 50) gd[0].data.shift();

          gd[0].data.push({
            x: cd.recorded_at,
            y: [cd.open, cd.high, cd.low, cd.close],
          });

          // if (vd[0].data.length > 50) vd[0].data.shift();

          vd[0].data.push({
            x: cd.recorded_at,
            y: cd.ask_volume,
          });
          vd[1].data.push({
            x: cd.recorded_at,
            y: cd.bid_volume,
          });

          setNeedNew(false);
        } else if (needNew() && cd.new_candle) {
          // if (gd[0].data.length > 50) gd[0].data.shift();

          gd[0].data.push({
            x: cd.recorded_at,
            y: [cd.open, cd.high, cd.low, cd.close],
          });

          // if (vd[0].data.length > 50) vd[0].data.shift();

          vd[0].data.push({
            x: cd.recorded_at,
            y: cd.ask_volume,
          });
          vd[1].data.push({
            x: cd.recorded_at,
            y: cd.bid_volume,
          });

          setNeedNew(true);
        } else if (!needNew() && !cd.new_candle) {
          gd[0].data[gd[0].data.length - 1] = {
            x: cd.recorded_at,
            y: [cd.open, cd.high, cd.low, cd.close],
          };

          vd[0].data[gd[0].data.length - 1] = {
            x: cd.recorded_at,
            y: cd.ask_volume,
          };
          vd[1].data[gd[0].data.length - 1] = {
            x: cd.recorded_at,
            y: cd.bid_volume,
          };
        } else if (!needNew() && cd.new_candle) {
          gd[0].data[gd[0].data.length - 1] = {
            x: cd.recorded_at,
            y: [cd.open, cd.high, cd.low, cd.close],
          };

          vd[0].data[gd[0].data.length - 1] = {
            x: cd.recorded_at,
            y: cd.ask_volume,
          };
          vd[1].data[gd[0].data.length - 1] = {
            x: cd.recorded_at,
            y: cd.bid_volume,
          };

          setNeedNew(true);
        }

        setGraphPriceData(gd);
        setGraphVolumeData(vd);
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
        <div class="flex flex-row gap-x-10">
          <SymbolInfoTable
            symbolInfo={symbolInfo}
            fetchErrorMsg={fetchErrorMsg}
          />
          <PlaceOrderForm />
        </div>
        <div class="flex flex-col gap-y-5 md:flex-row md:gap-x-5">
          {priceData.loading ? (
            <div>Loading...</div>
          ) : priceData.error ? (
            <div>Error: {priceData.error.message}</div>
          ) : (
            graphPriceData() && (
              <div class="flex flex-col gap-y-3 mb-10">
                <p>Price Chart</p>
                <SolidApexCharts
                  width={1000}
                  height={250}
                  type="candlestick"
                  options={{
                    chart: {
                      id: "price chart",
                      toolbar: {
                        show: false,
                      },
                    },
                    yaxis: {
                      tooltip: {
                        enabled: true,
                      },
                    },

                    xaxis: {
                      type: "datetime",
                    },
                  }}
                  series={graphPriceData() as ApexGraphPriceData}
                />
                <p>Volume Chart</p>
                <SolidApexCharts
                  width={1000}
                  height={250}
                  type="area"
                  options={{
                    chart: {
                      id: "volume chart",
                      animations: {
                        enabled: true,
                      },
                      toolbar: {
                        show: false,
                      },
                    },
                    xaxis: {
                      type: "datetime",
                    },
                    stroke: {
                      curve: "smooth",
                    },
                  }}
                  series={graphVolumeData() as ApexGraphVolumeData}
                />
                <UserPrivateInfoPanel />
              </div>
            )
          )}
        </div>
      </div>
    </div>
  );
};

export default Price;
