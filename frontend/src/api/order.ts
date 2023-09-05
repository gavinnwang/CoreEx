import { MsgResponse, sendPostRequest } from ".";
import { BASE_URL } from "../constants";

export type OrderType = "limit" | "market";
export type OrderSide = "buy" | "sell";

export type PlaceOrderParams = {
  price: number;
  volume: number;
  order_type: OrderType;
  order_side: OrderSide;
  symbol: string;
}

export async function placeOrder(
  params:  PlaceOrderParams,
  token: string
)  {
  const url = `http://${BASE_URL}/orders`;
  return sendPostRequest<MsgResponse>(url, params, token);
}

