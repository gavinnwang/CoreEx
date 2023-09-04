// store.js
import { createSignal } from "solid-js";
import { User } from "../api/user";


// const tokenLocal = localStorage.getItem("token");
export const [token, setToken] = createSignal< undefined| string >(undefined);
export const [user, setUser] = createSignal<null | User>(null);