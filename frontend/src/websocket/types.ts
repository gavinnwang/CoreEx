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

type PriceData = {
  recorded_at: number;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}[];

type GraphPriceDataPoint = {
  x: Date;
  y: [number, number, number, number];
}

type GraphPriceData = GraphPriceDataPoint[];

type SymbolInfo = {
  symbol: string;
  price: number;
  ask_volume: number;
  bid_volume: number;
  best_bid: number;
  best_ask: number;
};

type WSResponseGetMarketPrice = WSResponseBase & {
  result?: MarketPrice;
};

type WSResponseGetSymbolInfo = WSResponseBase & {
  result?: SymbolInfo;
};
