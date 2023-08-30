import { createSignal, type Component, createEffect } from "solid-js";

import logo from "./logo.svg";
import styles from "./App.module.css";

const App: Component = () => {
  const [price, setPrice] = createSignal(0);

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
    <div class={styles.App}>
      <header class={styles.header}>
        <img src={logo} class={styles.logo} alt="logo" />
        <p>market price: {price()}</p>
      </header>
    </div>
  );
};

export default App;
