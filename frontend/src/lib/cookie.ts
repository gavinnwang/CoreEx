import Cookies from "js-cookie";
import { COOKIE_NAME_JWT_TOKEN } from "../constants";

function setSessionCookie(token: string) {
  const expirationDate = new Date();
  expirationDate.setDate(expirationDate.getDate() + 30);

  Cookies.set(COOKIE_NAME_JWT_TOKEN, token, {
    path: "/",
    expires: expirationDate,
    httpOnly: false,
    secure: import.meta.env.VITE_ENV === "production",
  });
}

export { setSessionCookie };