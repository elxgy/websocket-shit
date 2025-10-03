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

  TEST_USERS: [
    { username: "yuri", password: "yuricbtt" },
    { username: "bernardo", password: "yuricbtt" },
    { username: "pedro", password: "yuricbtt" },
    { username: "marcelo", password: "yuricbtt" },
    { username: "giggio", password: "giggio123" },
    { username: "ramos", password: "ramosgay" },
    { username: "markin", password: "markinviado" },
  ],
};

export default config;
