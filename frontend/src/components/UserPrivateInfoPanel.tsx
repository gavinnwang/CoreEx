import { createResource } from "solid-js";

import Cookies from "js-cookie";
import { COOKIE_NAME_JWT_TOKEN } from "../constants";
import { getUserPrivateInfo } from "../api/user";
import { Refetch } from "./icons/Refetch";
import toast from "solid-toast";

function getToken(): string | undefined {
  const jwtToken = Cookies.get(COOKIE_NAME_JWT_TOKEN);
  return jwtToken;
}

const UserPrivateInfoPanel = () => {
  const [userInfo, { refetch }] = createResource(
    getToken(),
    getUserPrivateInfo
  );

  return (
    <div class="overflow-x-auto w-full h-fit border rounded-md shadow-sm">
      <div class="flex w-full justify-between py-1 font-semibold p-2">
        <div class="invisible">
          <Refetch />
        </div>
        Account
        <button
          onClick={() => {
            refetch();
            if (userInfo()?.cash_balance === 0) {
              toast.error("Log in to view account info");
              return;
            } else {
              toast.success("Refetched user info");
            }
          }}
        >
          <Refetch />
        </button>
      </div>

      <table class="table w-full text-sm">
        <tbody>
          <tr class="hover">
            <td class="font-semibold">Cash Balance</td>
            <td>${userInfo()?.cash_balance}</td>
          </tr>
        </tbody>
      </table>

      <div class="mt-4">
        <h3 class="font-semibold mx-4 my-2">Holdings</h3>
        <table class="table w-full text-sm">
          <thead>
            <tr>
              <th class="font-semibold">Symbol</th>
              <th class="font-semibold">Volume</th>
            </tr>
          </thead>
          <tbody>
            {userInfo()?.holdings ? (
              userInfo()?.holdings.map((holding) => (
                <tr class="hover">
                  <td>{holding.symbol}</td>
                  <td>{holding.volume}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td>No holdings</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div class="mt-4">
        <h3 class="font-semibold mx-4 my-2">Orders</h3>
        <table class="table w-full text-sm">
          <thead>
            <tr>
              {[
                "Symbol",
                "Order ID",
                "Side",
                "Status",
                "Type",
                "Price",
                "Avg. filled Price",
                "Filled Time",
                "Current Vol.",
                "Initial Vol.",
              ].map((header) => (
                <th class="font-semibold">{header}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {userInfo()?.orders ? (
              userInfo()?.orders.map((order) => (
                <tr class="hover">
                  <td>{order.symbol}</td>
                  <td>{order.order_id}</td>
                  <td>{order.order_side}</td>
                  <td>{order.order_status}</td>
                  <td>{order.order_type}</td>
                  <td>
                    ${order.order_type === "Market" ? "N/A" : order.price}
                  </td>
                  <td>
                    {order.filled_at
                      ? (order.total_processed / order.initial_volume).toFixed(2)
                      : "N/A"}
                  </td>

                  <td>{new Date(order.filled_at_time).toLocaleTimeString()}</td>
                  <td>{order.volume}</td>
                  <td>{order.initial_volume}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td>No orders</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default UserPrivateInfoPanel;
