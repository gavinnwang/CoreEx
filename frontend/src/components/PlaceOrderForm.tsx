import { createSignal } from "solid-js";
import {
  OrderSide,
  OrderType,
  PlaceOrderParams,
  placeOrder,
} from "../api/order";
import { token } from "../store";
import toast from "solid-toast";
import { APIError } from "../api";

export default function PlaceOrderForm() {
  const [order, setOrder] = createSignal<PlaceOrderParams>({
    price: 0,
    volume: 0,
    order_type: "limit",
    order_side: "sell",
    symbol: "",
  });

  const handleSubmit = (event: Event) => {
    event.preventDefault();
    const t = token();
    const o = order();
    if (!o.symbol) {
      toast.error("Please enter a symbol");
      return;
    }
    if (!o.price) {
      toast.error("Please enter a price");
      return;
    }
    if (!o.volume) {    
      toast.error("Please enter a volume");
      return;
    }
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
    <div class="w-content">
      <h1>Place an Order</h1>
      <form onSubmit={handleSubmit}>
        <label>
          Price:
          <input
            type="number"
            value={order().price}
            onInput={(e) =>
              setOrder({ ...order(), price: +e.currentTarget.value })
            }
          />
        </label>
        <label>
          Volume:
          <input
            type="number"
            value={order().volume}
            onInput={(e) =>
              setOrder({ ...order(), volume: +e.currentTarget.value })
            }
          />
        </label>
        <label>
          Order Type:
          <select
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
        </label>
        <label>
          Order Side:
          <select
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
        </label>
        <label>
          Symbol:
          <input
            type="text"
            value={order().symbol}
            onInput={(e) =>
              setOrder({ ...order(), symbol: e.currentTarget.value })
            }
          />
        </label>
        <button type="submit" class="btn">
          Submit Order
        </button>
      </form>
    </div>
  );
}
