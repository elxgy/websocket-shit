const express = require("express");
const path = require("path");

const app = express();
const PORT = process.env.PORT || 3000;

// Add request logging middleware
app.use((req, res, next) => {
  console.log(
    `${new Date().toISOString()} - ${req.method} ${req.url} from ${req.ip}`,
  );
  next();
});

// Serve static files from the React app build directory
app.use(
  express.static(path.join(__dirname, "build"), {
    maxAge: "1d", // Cache static files for 1 day
    etag: true,
  }),
);

// Health check endpoint for Railway
app.get("/health", (req, res) => {
  const healthData = {
    status: "ok",
    service: "frontend",
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    memory: process.memoryUsage(),
    port: PORT,
  };
  res.status(200).json(healthData);
});

// Error handling middleware
app.use((err, req, res, next) => {
  console.error("Server error:", err);
  res.status(500).json({
    error: "Internal server error",
    timestamp: new Date().toISOString(),
  });
});

// Catch all handler: send back React's index.html file for client-side routing
app.get("*", (req, res) => {
  try {
    res.sendFile(path.join(__dirname, "build", "index.html"));
  } catch (error) {
    console.error("Error serving index.html:", error);
    res.status(500).json({ error: "Failed to serve application" });
  }
});

// Graceful shutdown handling
process.on("SIGTERM", () => {
  console.log("SIGTERM received, shutting down gracefully");
  server.close(() => {
    console.log("Server closed");
    process.exit(0);
  });
});

process.on("SIGINT", () => {
  console.log("SIGINT received, shutting down gracefully");
  server.close(() => {
    console.log("Server closed");
    process.exit(0);
  });
});

const server = app.listen(PORT, "0.0.0.0", (err) => {
  if (err) {
    console.error("Failed to start server:", err);
    process.exit(1);
  }
  console.log(`🚀 Frontend server running on port ${PORT}`);
  console.log(`📡 Health check available at: http://localhost:${PORT}/health`);
  console.log(`🌍 Environment: ${process.env.NODE_ENV || "development"}`);
});
