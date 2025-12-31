import { useNavigate } from "react-router-dom";
import { Clock, Users } from "lucide-react";

const RecipeCard = ({ recipe }) => {
  const navigate = useNavigate();
  const recipeId = recipe.id || recipe.ID;
  
  // FIX: Explicitly parse and add times
 const displayTime = recipe.total_time || (Number(recipe.prep_time || 0) + Number(recipe.cook_time || 0));


  return (
    <div
      onClick={() => navigate(`/recipes/${recipeId}`)}
      className="cursor-pointer bg-white p-6 rounded-3xl shadow-sm hover:shadow-md transition-all border border-gray-100 flex flex-col h-full group"
    >
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-xl font-black text-slate-800 group-hover:text-orange-600 transition-colors">
            {recipe.name}
        </h3>
        <span className="bg-orange-50 text-orange-600 text-[10px] font-black uppercase px-2 py-1 rounded-lg">
            {recipe.category || "General"}
        </span>
      </div>

      <p className="text-lg text-slate-500 line-clamp-2 mb-6 flex-1 italic font-medium">
        {recipe.description || "No description provided."}
      </p>

      <div className="flex items-center justify-between text-slate-400 font-bold border-t pt-4">
        <div className="flex items-center gap-2 text-base">
            <Users size={18} />
            <span>{recipe.servings} ppl</span>
        </div>
        <div className="flex items-center gap-2 text-base">
            <Clock size={18} />
            {/* FIX: Shows calculated total time */}
            <span>{displayTime}m</span>
        </div>
      </div>
    </div>
  );
};

export default RecipeCard;