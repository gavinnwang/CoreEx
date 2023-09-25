import { Show, createSignal } from "solid-js";
import {
  OrderSide,
  OrderType,
  PlaceOrderParams,
  placeOrder,
} from "../api/order";
import { token } from "../store";
import toast from "solid-toast";
import { APIError } from "../api";


export default function PlaceOrderForm({ className }: { className?: string }) {
  const [order, setOrder] = createSignal<PlaceOrderParams>({
    price: 0,
    volume: 0,
    order_type: "market",
    order_side: "sell",
    symbol: "",
  });



  const handleSubmit = (event: Event) => {
    event.preventDefault();

    const o = order();
    if (!o.symbol) {
      toast.error("Please enter a symbol");
      return;
    }
    if (o.order_type === "limit" && !o.price) {
      toast.error("Please enter a price");
      return;
    }
    if (!o.volume) {
      toast.error("Please enter a volume");
      return;
    }

    const t = token();
    if (t) {
      console.log("Order submitted:", o);
      placeOrder(order(), t)
        .then((res) => {
          toast.success(res.message);
        })
        .catch((err) => {
          if (err instanceof APIError) {
            toast.error(err.message);
          }
        });
    } else {
      toast.error("Please login first");
    }
  };

  return (
    <div class={className}>
      <form class="flex flex-col gap-y-4" onSubmit={handleSubmit}>
        <div class="flex flex-row gap-x-2">
          <div class="flex flex-col w-48">
            <label>
              <span class="label-text font-bold">Order type</span>
            </label>
            <select
              class="select select-primary w-full max-w-xs"
              value={order().order_side}
              onChange={(e) =>
                setOrder({
                  ...order(),
                  order_side: e.currentTarget.value as OrderSide,
                })
              }
            >
              <option value="buy">Buy</option>
              <option value="sell">Sell</option>
            </select>
          </div>

          <label class="w-48">
            <span class="label-text font-bold">Volume</span>
            <input
              class="input input-bordered input-primary w-full max-w-xs"
              type="number"
              value={order().volume}
              onInput={(e) =>
                setOrder({ ...order(), volume: +e.currentTarget.value })
              }
            />
          </label>
          <label class="w-48">
            <span class="label-text font-bold">Symbol</span>
            <input
              class="input input-bordered input-primary w-full max-w-xs"
              type="text"
              value={order().symbol}
              onInput={(e) =>
                setOrder({ ...order(), symbol: e.currentTarget.value })
              }
            />
          </label>
        </div>

        <div class="flex flex-row gap-x-2">
          <div class="flex flex-col w-48">
            <label>
              <span class="label-text font-bold">Order type</span>
            </label>
            <select
              class="select select-primary w-full max-w-xs"
              value={order().order_type}
              onChange={(e) =>
                setOrder({
                  ...order(),
                  order_type: e.currentTarget.value as OrderType,
                })
              }
            >
              <option value="limit">Limit</option>
              <option value="market">Market</option>
            </select>
          </div>

          <Show when={order().order_type === "limit"}>
            <label class="w-48">
              <span class="label-text font-bold">Price</span>
              <input
                class="input input-bordered input-primary w-full max-w-xs"
                type="number"
                value={order().price}
                onInput={(e) =>
                  setOrder({ ...order(), price: +e.currentTarget.value })
                }
              />
            </label>
          </Show>
        </div>
        <button
          type="submit"
          disabled={
            !order().symbol ||
            !order().volume ||
            (order().order_type === "limit" && !order().price)
          }
          class="btn w-48 btn-secondary btn-outline"
        >
          Submit Order
        </button>
      </form>
    </div>
  );
}
