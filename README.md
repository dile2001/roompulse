RoomPulse is a simple IoT project for monitoring room environmental conditions such as temperature and humidity.

The project is being built incrementally to explore real-world IoT architecture, Go, device communication, monitoring, and cloud deployment.

Current Architecture

Sensor Simulator → JSON → HTTP → Go Server

The sensor simulator generates temperature readings and sends them to the RoomPulse server every few seconds.

Each reading contains:

Device ID
Temperature
Measurement timestamp

The server receives, decodes, validates, and processes the sensor readings.

Planned Evolution

Sensor Simulator
→ ESP32 + DHT22
→ Wi-Fi
→ HTTP / MQTT
→ Go Backend
→ Database
→ Monitoring & Dashboard
→ Docker / Kubernetes
→ Cloud

The goal is to gradually evolve RoomPulse from a small simulated sensor into a complete end-to-end IoT system.
