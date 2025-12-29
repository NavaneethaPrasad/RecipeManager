import { useEffect, useState } from "react";
import api from "../api/axios";
import RecipeCard from "../components/RecipeCard";
import { Link } from "react-router-dom";
import { Plus } from "lucide-react"; // Optional: Add icon back if you want

const Home = () => {
  const [recipes, setRecipes] = useState([]);

  useEffect(() => {
    api.get("/recipes")
      .then(res => setRecipes(res.data))
      .catch(err => console.error(err));
  }, []);

  return (
    <div className="p-8 bg-orange-50 min-h-screen">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-800">My Recipes</h1>

        {/* --- FIX IS HERE: Change 'to' path --- */}
        <Link
          to="/add-recipe" 
          className="bg-orange-600 text-white px-4 py-2 rounded flex items-center gap-2 hover:bg-orange-700 transition"
        >
          + Add Recipe
        </Link>
      </div>

      {recipes.length === 0 ? (
        <div className="text-center text-gray-500 mt-10">
            <p>No recipes found.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {recipes.map(recipe => (
            <RecipeCard key={recipe.id || recipe.ID} recipe={recipe} />
          ))}
        </div>
      )}
    </div>
  );
};

export default Home;