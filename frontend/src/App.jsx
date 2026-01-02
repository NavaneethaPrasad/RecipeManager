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
import ShoppingList from './pages/ShoppingList';
import ScaleRecipes from './pages/ScaleRecipes'; 

const LoadingSpinner = () => (
    <div className="flex h-screen items-center justify-center bg-orange-50">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-orange-600"></div>
    </div>
);

const ProtectedRoute = ({ children }) => {
    const { user, loading } = useAuth();
    const authHint = localStorage.getItem('auth_hint');

    if (loading) {
        return <LoadingSpinner />;
    }

    if (!user && !authHint) {
        return <Navigate to="/login" replace />;
    }

    if (!user) {
        return <Navigate to="/login" replace />;
    }

    return children;
};

const PublicRoute = ({ children }) => {
    const { user, loading } = useAuth();
    
    if (loading) return <LoadingSpinner />;
    
    if (user) return <Navigate to="/" replace />; 
    
    return children;
};

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
            {/* Public Routes */}
            <Route path="/login" element={<PublicRoute><Login /></PublicRoute>} />
            <Route path="/register" element={<PublicRoute><Register /></PublicRoute>} />
            
            {/* Protected Routes*/}
            <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>} />
            <Route path="/recipes" element={<ProtectedRoute><RecipeList /></ProtectedRoute>} />
            <Route path="/add-recipe" element={<ProtectedRoute><AddRecipe /></ProtectedRoute>} />
            <Route path="/recipes/:id" element={<ProtectedRoute><RecipeDetails /></ProtectedRoute>} />
            <Route path="/meal-planner" element={<ProtectedRoute><MealPlanner /></ProtectedRoute>} />
            <Route path="/shopping-list" element={<ProtectedRoute><ShoppingList /></ProtectedRoute>} />
            <Route path="/scale" element={<ProtectedRoute><ScaleRecipes /></ProtectedRoute>} />
            
            <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
        <Toaster position="top-right" />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;