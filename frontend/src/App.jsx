import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { Toaster } from 'react-hot-toast';

import Login from './pages/Login';
import Register from './pages/Register';
import Home from './pages/Home';
import RecipeList from './pages/RecipeList';
import AddRecipe from './pages/AddRecipe';
import RecipeDetails from './pages/RecipeDetails';
import MealPlanner from './pages/MealPlanner';
import ShoppingList from './pages/ShoppingList'; // ✅ Import Real Page

// Placeholder for Scale (We will do this next)
const ScaleRecipes = () => <div className="p-10 text-center">Scale Recipes Coming Soon...</div>;

const ProtectedRoute = ({ children }) => {
    const { user, loading } = useAuth();
    if (loading) return <div>Loading...</div>;
    if (!user) return <Navigate to="/login" />;
    return children;
};

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            
            <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
            <Route path="/recipes" element={<ProtectedRoute><RecipeList /></ProtectedRoute>} />
            <Route path="/add-recipe" element={<ProtectedRoute><AddRecipe /></ProtectedRoute>} />
            <Route path="/recipes/:id" element={<ProtectedRoute><RecipeDetails /></ProtectedRoute>} />
            <Route path="/meal-planner" element={<ProtectedRoute><MealPlanner /></ProtectedRoute>} />
            
            {/* ✅ Update Route */}
            <Route path="/shopping-list" element={<ProtectedRoute><ShoppingList /></ProtectedRoute>} />
            
            <Route path="/scale" element={<ProtectedRoute><ScaleRecipes /></ProtectedRoute>} />
        </Routes>
        <Toaster position="top-right" />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;