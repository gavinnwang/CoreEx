import { sendGetRequest } from ".";
import { BASE_URL } from "../constants";

export async function getPriceData(
    symbol: string
) : Promise<SymbolInfo[]> {
  const url = `http://${BASE_URL}/price-history/${symbol}`;
  return sendGetRequest<SymbolInfo[]>(url);
}

