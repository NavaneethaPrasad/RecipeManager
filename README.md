# MealMate: Your Intelligent Kitchen Assistant

**MealMate** is a full-stack web application designed to streamline recipe management, weekly meal planning, and automated grocery shopping. Built with a focus on **Data Integrity**, **Security**, and **Responsive Design**, it allows users to organize their culinary life with mathematical precision.

## 🌟 Core Features

### 1. Advanced Recipe Management (CRUD)
* **Dynamic Creation:** Create and edit recipes with dynamic ingredient lists and step-by-step instructions.
* **Smart Data Mapping:** Categorize recipes and track preparation/cooking times.
* **Safe Deletion:** Implemented database transactions to handle cascading deletes of linked instructions and meal plans to ensure data integrity.

### 2. Interactive Weekly Meal Planner
* **Weekly Grid:** A 7-day responsive grid to assign recipes to specific meal slots (Breakfast, Lunch, Dinner, Snack).
* **Targeted Servings:** Set specific serving sizes for each meal, independent of the base recipe.
* **Horizontal Swipe:** Mobile-optimized calendar view for seamless planning on the go.

### 3. Automated Shopping List Generation
* **Aggregation Logic:** The system scans your weekly plan and combines identical ingredients into a single line item.
* **The Scaling Engine:** Automatically calculates ingredient quantities using the ratio:  
  ` (Meal Plan Target Servings / Recipe Base Servings) * Ingredient Amount `
* **Interactive Checklist:** Track your progress in the store with persistent checkboxes.

### 4. Real-time Ingredient Scaling
* **Instant Preview:** Preview ingredient adjustments instantly without making database calls.
* **Versatility:** Ideal for adjusting portions for large gatherings or solo meals.

---

## 🛠 Tech Stack

### Backend
* **Language:** Golang (Go)
* **Framework:** Gin Gonic
* **ORM:** GORM
* **Security:** JWT Authentication with **httpOnly Cookies** (Protection against XSS)

### Frontend
* **Framework:** React (Vite)
* **Styling:** Tailwind CSS 

### DevOps & Tools
* **Database:** PostgreSQL
* **Containerization:** Docker & Docker Compose

---

## ⚙️ Installation & Setup

### Prerequisites
* Docker & Docker Compose
* Git

### 1. Clone the Repository
```bash
git clone git@github.com:NavaneethaPrasad/RecipeManager.git
cd RecipeManager
```

### 2. Configure Environment

**Root Directory Config:** Create a `.env` file in the root directory (`RecipeManager/.env`) and paste the following:
```env
DB_NAME=recipe_db
DB_USER=recipe_user
DB_PASS=recipe_pass
```

**Backend Config:** Navigate to the backend folder and create a `.env` file (`RecipeManager/backend/.env`):
```env
DB_HOST=postgres 
DB_USER=recipe_user
DB_PASSWORD=recipe_pass
DB_NAME=recipe_db
DB_PORT=5432
PORT=8080
JWT_SECRET=supersecretkey
```

### 3. Build & Run
Run the following command to build images and start the containers:
```bash
docker compose up --build -d
```
> **Note:** The first run may take a few minutes to download the Go and Node images.

### 4. Access the App
Once the logs show the servers are running, open your browser:
* **Frontend UI:** `http://localhost`

---

## 📖 User Manual: How to use MealMate

### 1. Account Setup
* **Registration:** Click on "Sign Up" from the login page. Enter your name, email, and a secure password (minimum 6 characters).
* **Login:** Use your credentials to access your personal dashboard. The system uses secure cookies, so you will stay logged in even if you refresh the page.
* **Logout:** Use the Logout button in the top navigation bar to securely end your session and clear your data from the browser.

### 2. Managing Your Recipes
* **Add a Recipe:** Click "Add Recipe" from the Dashboard or the Recipes page.
* **Ingredients:** Use the "+ Add Ingredient" button to create as many rows as you need. Enter the item name, quantity, and unit (e.g., "Onion", "2", "pcs").
* **Instructions:** Type your steps in the Instructions box. **Tip:** Press Enter to start a new line for each step; the app will automatically number them for you.
* **View & Edit:** Click on any recipe card to see the full details. Use the Edit (Pencil icon) to update ingredients or the Delete (Trash icon) to remove it.

### 3. Planning Your Weekly Meals
* **The Grid:** Go to the Meal Planner. You will see a 7-day calendar view.
* **Adding Meals:** Click the "+" button in any slot (Breakfast, Lunch, Dinner, Snack). Select a recipe from your collection.
* **Crucial Step:** Enter the **Desired Servings**. If you are cooking for 4 people today, enter "4". The shopping list will use this number to calculate ingredients.
* **Navigation:** Use the Left/Right arrows at the top to plan for next week or look at previous plans.

### 4. Generating a Shopping List
* **Generate:** Go to the Shopping List page and select your date range.
* **One-Click List:** Click "Generate". MealMate will scan your planned meals, scale the ingredients for your target servings, and combine duplicates.
* **Checklist:** As you shop, tap the ingredient name to check it off.

### 5. Scaling Portions (Quick Tool)
* **Instant Math:** Use the Scale tool to adjust a recipe's portions instantly without adding it to your planner.
* **Result:** The table updates in real-time to show the Original Amount vs. the New Scaled Amount.

---

## 🧪 Running Tests
To run the backend integration tests (requires Go installed locally):

1. Ensure the Docker database container is running.
2. Navigate to the backend folder: `cd backend`
3. Run the tests:
```bash
go test -v ./internal/services
```

---

## 🛠 Troubleshooting (Docker Edition)

### 1. Check if Containers are Running
```bash
docker-compose ps
```

### 2. View Service Logs
* **Backend:** `docker-compose logs recipe_backend`
* **Database:** `docker-compose logs recipe_postgres`
* **Frontend:** `docker-compose logs recipe_frontend`

### 3. Common Issues & Solutions
* **"Database Connection Refused":** The Go backend started before PostgreSQL was ready.  
  *Solution:* `docker-compose restart backend`
* **Port Conflict:** Another instance is using 8080 or 5173.  
  *Solution:* `docker-compose down` and then `docker-compose up -d`
* **Changes not reflecting:** Docker is using an old image.  
  *Solution:* `docker-compose up -d --build`
* **Frontend cannot talk to Backend:** baseURL in `axios.js` might be wrong.  
  *Solution:* Ensure baseURL is `http://localhost:8080/api`.

### 4. Resetting the Environment
If you want to clear all data and start fresh:
```bash
docker-compose down -v
docker-compose up -d
```
```