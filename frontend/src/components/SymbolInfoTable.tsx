import { Accessor } from "solid-js";

export default function SymbolInfoTable({
  symbolInfo,
  fetchErrorMsg,
  className,
}: {
  symbolInfo: Accessor<SymbolInfo | null>;
  fetchErrorMsg: Accessor<string | null>;
  className?: string;
}) {
  return (
    <div class="overflow-x-auto w-80 h-fit border rounded-md shadow-sm ">
      <div class="flex w-full items-center justify-center pt-1 font-semibold h-8">
       {symbolInfo()?.symbol}
      </div>
      <table class="table">
        <tbody>
          <tr class="hover">
            <td class="text-sm font-semibold">Market price</td>
            <td>${symbolInfo()?.price}</td>
          </tr>
          <tr class="hover">
            <td class="text-sm font-semibold">Ask volume</td>
            <td>{symbolInfo()?.ask_volume}</td>
          </tr>

          <tr class="hover">
            <td class="text-sm font-semibold">Bid volume</td>
            <td>{symbolInfo()?.bid_volume}</td>
          </tr>

          <tr class="hover">
            <td class="text-sm font-semibold">Best Bid</td>
            <td>${symbolInfo()?.best_bid}</td>
          </tr>

          <tr class="hover">
            <td class="text-sm font-semibold">Best Ask</td>
            <td>${symbolInfo()?.best_ask}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
