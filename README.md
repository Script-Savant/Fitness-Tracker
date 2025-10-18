🏋️‍♂️ Fitness Tracker (Go + Gin + GORM)

Fitness Tracker is a lightweight web application built with Go, Gin, and GORM, designed to help users track workouts, monitor fitness metrics (like weight, height, and BMI), and manage personal fitness data securely.

It features user authentication, BMI calculation, and a clean Bootstrap-powered UI for smooth user experience.

📁 Project Structure
.
├── admin/                     # (Reserved for admin-related logic)
├── cmd/admin/main.go          # Alternative app entry point for admin panel
├── config/
│   └── database.go            # GORM database setup
├── fitness_tracker.db         # SQLite database (auto-generated)
├── handlers/                  # Controller logic
│   ├── auth.go                # Handles user registration/login
│   ├── metrics.go             # Handles BMI metrics creation/display
│   └── workout.go             # Handles workout CRUD operations
├── main.go                    # Main application entry point
├── middleware/
│   └── auth.go                # Authentication middleware (login protection)
├── models/
│   └── models.go              # GORM models: User, Workout, Metrics
├── routes/
│   ├── auth.go                # Authentication routes
│   ├── metrics.go             # Metrics routes
│   └── workouts.go            # Workout routes
├── static/                    # Static assets
│   ├── CSS/
│   │   └── bootstrap.min.css  # Bootstrap styling
│   └── Images/
│       └── icon.png           # App icon
├── templates/
│   ├── layout.html            # Shared layout (header, footer, etc.)
│   ├── home.html              # Home dashboard
│   ├── auth/                  # Authentication pages
│   │   ├── login.html
│   │   └── register.html
│   ├── metrics/               # Metrics pages
│   │   ├── create-metrics.html
│   │   └── display-metrics.html
│   └── workout/               # Workout pages
│       ├── create.html
│       └── update.html
├── utils/
│   ├── auth.go                # Session and user utilities
│   └── templates.go           # Multi-template renderer setup
├── routes/
│   ├── auth.go
│   ├── metrics.go
│   └── workouts.go
├── structure.pdf              # System design/architecture diagram
├── structure_page-0001.jpg    # Project structure snapshot
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
└── README.md                  # You are here 🚀


⚙️ Features

✅ User Authentication

Register and log in with secure password hashing

Session-based user management

✅ Workouts Management

Create and update workout sessions

Track workout type, duration, distance, and completion status

✅ Fitness Metrics Tracking

Record weight and height

Automatically calculate BMI (Body Mass Index)

Automatically classify BMI status (Underweight, Normal, Overweight, Obese)

✅ Modern UI

Built with Bootstrap 5 for responsive design

Clean card-based layout and intuitive navigation

✅ Persistence

Data stored in SQLite using GORM ORM

✅ Template Rendering

Modular templates with layout.html

multitemplate.Renderer for route-specific templates

Custom template functions for better display (add, etc.)

🧠 Tech Stack
| Layer             | Technology                              |
| ----------------- | --------------------------------------- |
| Backend Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM               | [GORM](https://gorm.io/)                |
| Database          | SQLite (default)                        |
| Frontend          | HTML5, Bootstrap 5                      |
| Templates         | Go `html/template` with `multitemplate` |
| Auth              | Custom session-based middleware         |
| Language          | Go 1.22+                                |


🚀 Getting Started

1️⃣ Clone the Repository
```git clone https://github.com/Script-Savant/Fitness-Tracker.git```
```cd Fitness-Tracker```

2️⃣ Install Dependencies
```go mod tidy```

3️⃣ Run the Application
```go run main.go```

4️⃣ Access the App

Open your browser and go to:
```http://localhost:8080```


🧩 Main Routes
| Route                 | Method   | Description          |
| --------------------- | -------- | -------------------- |
| `/register`           | GET/POST | User registration    |
| `/login`              | GET/POST | User login           |
| `/logout`             | GET      | Log out current user |
| `/home`               | GET      | User dashboard       |
| `/create-workout`     | GET/POST | Add a workout        |
| `/update-workout/:id` | GET/POST | Edit a workout       |
| `/create-metrics`     | GET/POST | Record new BMI data  |
| `/display-metrics`    | GET      | View all metrics     |

🧑‍💻 Developer Notes

Make sure your templates are registered in utils.SetupTemplates().

Static files (Bootstrap, images) are served under /static.

Middleware in middleware/auth.go protects private routes.

The utils/auth.go file manages session and user retrieval logic.

🧾 License

MIT License © 2025 Alex Kinuthia


🧮 BMI Calculation Logic

In the Metrics model (models/models.go):
```func (m *Metrics) BeforeSave(tx *gorm.DB) (err error) {
	if m.HeightCm <= 0 {
		m.BMI = 0
		m.Status = "Invalid height"
		return nil
	}

	divisor := math.Pow(float64(m.HeightCm)/100, 2)
	m.BMI = float32(float64(m.WeightKg) / divisor)

	switch {
	case m.BMI < 18.50:
		m.Status = "Underweight"
	case m.BMI >= 18.50 && m.BMI < 25.00:
		m.Status = "Normal"
	case m.BMI >= 25.00 && m.BMI < 30.00:
		m.Status = "Overweight"
	default:
		m.Status = "Obese"
	}
	return nil
}```




`Feel Free to add more functionality`