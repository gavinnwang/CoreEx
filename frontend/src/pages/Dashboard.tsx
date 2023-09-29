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

  const [fetchErrorMsg, setFetchErrorMsg] = createSignal<string | null>(null);

  const [graphPriceData, setGraphPriceData] =
    createSignal<ApexGraphPriceData | null>(null);

  const [graphVolumeData, setGraphVolumeData] =
    createSignal<ApexGraphVolumeData | null>(null);

  const symbol = "AAPL";

  const [priceData] = createResource(symbol, getPriceHistoryData);

  const convertDataToPriceGraphFormat = (
    data: CandleDataPoint[]
  ): { data: GraphPriceDataPoint[] }[] => {
    return [
      {
        data: data.map((item) => ({
          x: new Date(item.recorded_at * 1000),
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

  const convertDataToVolumeGraphFormat = (
    data: CandleDataPoint[]
  ): { data: GraphVolumeDataPoint[] }[] => {
    return [
      {
        data: data.map((item) => ({
          x: new Date(item.recorded_at * 1000),
          y: item.volume,
        })),
      },
    ];
  };

  async function getPriceHistoryData(
    symbol: string
  ): Promise<ApexGraphPriceData> {
    try {
      const data = await getPriceData(symbol);
      const graphPriceData = convertDataToPriceGraphFormat(data);
      const graphVolumeData = convertDataToVolumeGraphFormat(data);
      setGraphVolumeData(graphVolumeData);
      setGraphPriceData(graphPriceData);
      return graphPriceData;
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
      if (resData.success && resData.result) {
        setSymbolInfo(resData.result);
        // setCandleData(resData.result.candle_data);
        const cd = resData.result.candle_data;
        if (cd) {
          const gd = [...(graphPriceData() as ApexGraphPriceData)];
          const vd = [...(graphVolumeData() as ApexGraphVolumeData)];

          if (needNew() && !cd.new_candle) {
            if (gd[0].data.length > 50) gd[0].data.shift();

            gd[0].data.push({
              x: new Date(cd.recorded_at * 1000),
              y: [cd.open, cd.high, cd.low, cd.close],
            });

            if (vd[0].data.length > 50) vd[0].data.shift();

            vd[0].data.push({
              x: new Date(cd.recorded_at * 1000),
              y: cd.volume,
            });

            setNeedNew(false);
          }

          if (needNew() && cd.new_candle) {
            if (gd[0].data.length > 50) gd[0].data.shift();

            gd[0].data.push({
              x: new Date(cd.recorded_at * 1000),
              y: [cd.open, cd.high, cd.low, cd.close],
            });

            if (vd[0].data.length > 50) vd[0].data.shift();

            vd[0].data.push({
              x: new Date(cd.recorded_at * 1000),
              y: cd.volume,
            });

            setNeedNew(true);
          }

          if (!needNew() && !cd.new_candle) {
            gd[0].data[gd[0].data.length - 1] = {
              x: new Date(cd.recorded_at * 1000),
              y: [cd.open, cd.high, cd.low, cd.close],
            };

            vd[0].data[gd[0].data.length - 1] = {
              x: new Date(cd.recorded_at * 1000),
              y: cd.volume,
            };
          }

          if (!needNew() && cd.new_candle) {
            gd[0].data[gd[0].data.length - 1] = {
              x: new Date(cd.recorded_at * 1000),
              y: [cd.open, cd.high, cd.low, cd.close],
            };

            vd[0].data[gd[0].data.length - 1] = {
              x: new Date(cd.recorded_at * 1000),
              y: cd.volume,
            };

            setNeedNew(true);
          }

          setGraphPriceData(gd);
          setGraphVolumeData(vd);
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
          ) : (
            graphPriceData() && (
              <div class="flex flex-col gap-y-3">
                <p>Price Chart</p>
                <SolidApexCharts
                  width="800"
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
                  }}
                  series={graphPriceData() as ApexGraphPriceData}
                />
                <p>Volume Chart</p>
                <SolidApexCharts
                  width="800"
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
                    stroke: {
                      curve: "smooth",
                    },
                  
                  }}
                  series={graphVolumeData() as ApexGraphVolumeData}
                />
              </div>
            )
          )}
        </div>
      </div>
    </div>
  );
};

export default Price;
