import { useNavigate } from "react-router-dom";
import { Clock, Users } from "lucide-react";

const RecipeCard = ({ recipe }) => {
  const navigate = useNavigate();
  const recipeId = recipe.id || recipe.ID;

  // No need to recalculate total time here, backend sends it as "total_time"
  const totalTime = recipe.total_time || (recipe.prep_time + recipe.cook_time);

  return (
    <div
      onClick={() => navigate(`/recipes/${recipeId}`)}
      className="cursor-pointer bg-white p-5 rounded-xl shadow-md hover:shadow-xl transition border border-gray-100 flex flex-col h-full"
    >
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-lg font-bold text-gray-800 line-clamp-1">{recipe.name}</h3>
        
        {/* FIX: Use the 'category' from the API response */}
        <span className="text-xs font-semibold bg-orange-100 text-orange-600 px-2 py-1 rounded-full">
            {recipe.category || "General"}
        </span>
      </div>

      <p className="text-gray-500 text-sm line-clamp-2 mb-4 flex-1">
        {recipe.description || "No description provided."}
      </p>

      <div className="mt-auto flex items-center justify-between text-sm text-gray-500 border-t pt-3">
        <div className="flex items-center gap-1">
            <Users size={16} />
            <span>{recipe.servings} ppl</span>
        </div>
        <div className="flex items-center gap-1">
            <Clock size={16} />
            {/* FIX: Display time correctly, show "N/A" only if zero */}
            <span>{totalTime > 0 ? `${totalTime} min` : "N/A"}</span>
        </div>
      </div>
    </div>
  );
};

export default RecipeCard;