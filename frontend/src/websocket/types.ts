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


type CandleDataPoint = {
  recorded_at: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}

type GraphPriceDataPoint = {
  x: string;
  y: [number, number, number, number];
};

type GraphVolumeDataPoint = {
  x: string;
  y: number;
}


type ApexGraphPriceData = {
  data: GraphPriceDataPoint[];
}[];

type ApexGraphVolumeData = {
  name: string;
  data: GraphVolumeDataPoint[];
}[];

type SymbolInfo = CandleDataUpdate & {
  symbol: string;
  price: number;
  ask_volume: number;
  bid_volume: number;
  best_bid: number;
  best_ask: number;
  
};

type CandleDataUpdate = CandleDataPoint & {
  new_candle: boolean;
};

type WSResponseGetMarketPrice = WSResponseBase & {
  result?: MarketPrice;
};

type WSResponseGetSymbolInfo = WSResponseBase & {
  result?: SymbolInfo;
};

// {
//   "cash_balance": 100025,
//   "holdings": [
//       {
//           "symbol": "AAPL",
//           "volume": "-0.5"
//       }
//   ],
//   "orders": [
//       {
//           "symbol": "AAPL",
//           "order_id": "01HBG15BFQ79TVJRRQ7HE32A2V",
//           "order_side": "Sell",
//           "order_status": "Filled",
//           "order_type": "Limit",
//           "filled_at": 50,
//           "total_processed": 25,
//           "volume": 0,
//           "initial_volume": 0.5,
//           "price": 50
//       }
//   ]
// }

type UserPrivateInfo = {
  cash_balance: number;
  holdings: {
    symbol: string;
    volume: number;
  }[];
  orders: {
    symbol: string;
    order_id: string;
    order_side: string;
    order_status: string;
    order_type: string;
    filled_at: number;
    filled_at_time: string;
    total_processed: number;
    volume: number;
    initial_volume: number;
    price: number;
  }[];
}