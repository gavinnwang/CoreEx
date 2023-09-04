// import { createSignal } from "solid-js";
// import { OAuth2Client } from "google-auth-library";

// function GoogleLogin() {
//   const [loggedIn, setLoggedIn] = createSignal(false);
//   const [userInfo, setUserInfo] = createSignal(null);

//   const handleLogin = async () => {
//     const client = new OAuth2Client(
//     );

//     // You'll likely use a popup or redirect flow here
//     const url = client.generateAuthUrl({
//       access_type: "offline",
//       scope: "https://www.googleapis.com/auth/userinfo.profile",
//     });

//     // Redirect or open a popup for login
//     window.location.href = url;
//   };
//   return (
//     <div>
//       <button onClick={handleLogin}>Login with Google</button>
//     </div>
//   )
// }

// export default GoogleLogin;
