type WSResponseBase = {
  event: string;
  success: boolean;
  error_message?: string;
};

type WSResponse = {
  event: string;
};

type ParamsStreamPrice = WSResponse & {
  params: {
    symbol: string;
  };
};

type MarketPrice = {
  symbol: string;
  price: number;
};

type SymbolInfo = {
  symbol: string;
  price: number;
  ask_volume: number;
  bid_volume: number;
  best_bid: number;
  best_ask: number;
}

type WSResponseGetMarketPrice = WSResponseBase & {
  result?: MarketPrice;
};

type WSResponseGetSymbolInfo = WSResponseBase & {
  result?: SymbolInfo
}
