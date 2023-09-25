import { SolidApexCharts } from "solid-apexcharts";
import { createSignal } from "solid-js";
export default function CandleGraph({
  data,
}: {
  data: ApexAxisChartSeries;
}) {
  const [options] = createSignal({
    chart: {
      id: "solidchart-example",
    },
  });

    console.log(data);
    return (
      <SolidApexCharts
        width="800"
        type="candlestick"
        options={options()}
        series={data}
      />
    );

}
