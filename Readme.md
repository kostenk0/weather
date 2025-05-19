# 🌤️ Weather Subscription API

This application allows users to subscribe to regular weather updates for a selected city and view the current forecast.

🔗 **Live demo:** [https://weather-phz6.onrender.com/](https://weather-phz6.onrender.com/)

---

## 🔧 Features

- 🔍 **Get current weather** by city via `/api/weather?city=...`
- 📩 **Subscribe to weather updates** via email with selected frequency (hourly or daily)
- ✅ **Confirm subscription** via email link
- ❌ **Unsubscribe** with a single click
- 🗃️ **Weather is cached** in PostgreSQL and automatically updated
- 📬 Emails are sent through a configured SMTP server
- 🐳 Fully containerized via Docker and deployable to Render

---

## 🔗 API Endpoints

| Method | Path                       | Description                    |
|--------|----------------------------|--------------------------------|
| GET    | `/api/weather?city=Kyiv`   | Get current weather for a city |
| POST   | `/api/subscribe`           | Subscribe to weather updates   |
| GET    | `/api/confirm/{token}`     | Confirm email subscription     |
| GET    | `/api/unsubscribe/{token}` | Unsubscribe from updates       |

---

## ⚙️ Local Setup

### 1. Clone the repository

```bash
git clone https://github.com/kostenk0/weather
cd weather-api
```

### 2. Create a .env file

<pre lang="dotenv"><code>
# .env.example

PORT=8080
APP_BASE_URL=http://localhost:8080

# SMTP (e.g. Gmail SMTP + App Password)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your@email.com
SMTP_PASS=your_app_password
EMAIL_FROM=Weather Service &lt;your@email.com&gt;

# Weather API
WEATHER_API_KEY=your_weatherapi_com_key
WEATHER_API_URL=https://api.weatherapi.com/v1

# PostgreSQL connection string
DATABASE_URL=postgres://weather_user:password@localhost:5432/weatherdb?sslmode=disable
</code></pre>

### 3. Run with Docker

```bash
docker-compose up --build
```

### 4. Test the app

- Open in browser: [http://localhost:8080](http://localhost:8080)

- Or test via curl:

```bash
curl http://localhost:8080/api/weather?city=Kyiv
```
