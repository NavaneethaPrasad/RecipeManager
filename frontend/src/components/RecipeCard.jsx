import { useNavigate } from "react-router-dom";
import { Clock, Users } from "lucide-react";

const RecipeCard = ({ recipe }) => {
  const navigate = useNavigate();
  const recipeId = recipe.id || recipe.ID;
  const displayTime = recipe.total_time || (Number(recipe.prep_time || 0) + Number(recipe.cook_time || 0));

  return (
    <div
      onClick={() => navigate(`/recipes/${recipeId}`)}
      className="cursor-pointer bg-white p-5 md:p-6 rounded-[2rem] shadow-sm hover:shadow-md active:scale-[0.97] transition-all border border-slate-100 flex flex-col h-full group"
    >
      <div className="flex justify-between items-start gap-2 mb-3">
        <h3 className="text-lg md:text-xl font-black text-slate-800 group-hover:text-orange-600 transition-colors leading-tight line-clamp-1">
            {recipe.name}
        </h3>
        <span className="flex-shrink-0 bg-orange-50 text-orange-600 text-[9px] md:text-[10px] font-black uppercase px-2 py-1 rounded-lg tracking-wider">
            {recipe.category || "General"}
        </span>
      </div>

      <p className="text-sm md:text-lg text-slate-500 line-clamp-2 mb-6 flex-1 italic font-medium leading-relaxed">
        {recipe.description || "No description provided."}
      </p>
      <div className="flex items-center justify-between text-slate-400 font-bold border-t border-slate-50 pt-4 mt-auto">
        <div className="flex items-center gap-1.5 md:gap-2 text-xs md:text-base">
            <Users size={16} className="md:w-[18px] md:h-[18px]" />
            <span>{recipe.servings} ppl</span>
        </div>
        <div className="flex items-center gap-1.5 md:gap-2 text-xs md:text-base">
            <Clock size={16} className="md:w-[18px] md:h-[18px]" />
            <span>{displayTime}m</span>
        </div>
      </div>
    </div>
  );
};

export default RecipeCard;