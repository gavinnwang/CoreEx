import { createSignal, type Component, createEffect } from "solid-js";


const Price: Component = () => {
  const [price, setPrice] = createSignal<number | null>(null);

  createEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/price");
    // Set up event listeners
    ws.addEventListener("open", () => {
      console.log("WebSocket connection opened");
    });

    ws.addEventListener("message", (event) => {
      const receivedMessage = event.data;
      console.log("Received:", receivedMessage);

      // Update state with the latest message
      setPrice(receivedMessage);
    });

    // Clean up: Close the WebSocket connection when this effect is destroyed
    return () => {
      console.log("WebSocket connection closed");
      ws.close();
    };
  });
  return (
    <div class="flex bg-sky-700 justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      <header class="flex items-center flex-col">
        {/* <img src={logo} class="w-40 h-40" alt="logo" /> */}
        <p class="text-sky-200 italic underline-offset-4">
          market price: {price() ?? "price data not available"}
        </p>
      </header>
    </div>
  );
};

export default Price;
