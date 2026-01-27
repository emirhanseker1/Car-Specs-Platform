# Car Specs Platform

A comprehensive platform for viewing detailed technical specifications of vehicles, built with a modern Go backend and a React frontend.

![Car Specs Screenshot](frontend/public/hero/hero-2.jpg)

## Features

-   **Detailed Vehicle Database**: Specifications including engine, transmission, dimensions, and more.
-   **Advanced Search**: Filter vehicles by brand, model, year, fuel type, and transmission.
-   **Comparison Tool**: Side-by-side comparison of different vehicle trims.
-   **Immersive UI**: Modern dark mode design with glassmorphism effects and high-quality imagery.
-   **Data Synchronization**: Automated sync with CarQuery API for up-to-date specs.

## Tech Stack

-   **Backend**: Go (Golang), SQLite (GORM), Chi Router
-   **Frontend**: React, TypeScript, Tailwind CSS, Framer Motion
-   **Tools**: Vite, Lucide React

## Getting Started

### Prerequisites

-   Go 1.22+
-   Node.js 18+ & npm

### One-Click Start (Recommended)

Simply double-click the `start.bat` file in the root directory. This will verify dependencies and launch both the backend and frontend services automatically.

### Manual Setup

#### 1. Backend

The backend handles the database and data synchronization.

```bash
cd backend
go mod tidy

# Run the server (scrapes data on first run if configured)
go run cmd/api/main.go
# Or to force a scrape: go run cmd/api/main.go -scrape
```

Server runs on: `http://localhost:8080`

#### 2. Frontend

The frontend is the user interface.

```bash
cd frontend
npm install
npm run dev
```

App runs on: `http://localhost:5173`

## Project Structure

-   `/backend`: Go API server and database logic.
-   `/frontend`: React application.
-   `/data`: Shared data resources (if any).

## API Endpoints

-   `GET /api/vehicles`: List brands and initial vehicle data.
-   `GET /api/trims/{id}`: detailed specs for a specific trim.
-   `GET /api/search`: Advanced search with filters.
-   `GET /api/featured`: Featured vehicles for homepage.

## License

This project is licensed under the MIT License.
