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

type ResultGetMarketPrice = {
  symbol: string;
  price: number;
};

type ResponseGetMarketPrice = WSResponseBase & {
  result?: ResultGetMarketPrice;
};
