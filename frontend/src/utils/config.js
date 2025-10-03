const config = {
  API_URL:
    process.env.REACT_APP_API_URL ||
    (process.env.NODE_ENV === "production"
      ? "https://websocket-shit-production.up.railway.app"
      : "http://localhost:8080"),

  WS_URL:
    process.env.REACT_APP_WS_URL ||
    (process.env.NODE_ENV === "production"
      ? "wss://websocket-shit-production.up.railway.app"
      : "ws://localhost:8080"),

  RECONNECT_ATTEMPTS: 5,
  RECONNECT_DELAY: 3000,

  MAX_MESSAGE_LENGTH: 500,
  MESSAGE_HISTORY_LIMIT: 100,
};

export default config;
