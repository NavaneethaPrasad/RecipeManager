import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import api from "../api/axios";
import toast from "react-hot-toast";
import { Edit, Trash2, Clock, Users, ArrowLeft } from "lucide-react";

const RecipeDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [recipe, setRecipe] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get(`/recipes/${id}`)
      .then(res => {
        setRecipe(res.data);
        setLoading(false);
      })
      .catch((err) => {
        console.error(err);
        toast.error("Recipe not found");
        navigate("/");
      });
  }, [id, navigate]);

  const handleDelete = async () => {
    if (!window.confirm("Delete this recipe permanently?")) return;

    try {
      await api.delete(`/recipes/${id}`);
      toast.success("Recipe deleted");
      navigate("/");
    } catch (error) {
      toast.error("Delete failed");
    }
  };

  if (loading || !recipe) return <div className="p-10 text-center">Loading...</div>;

  // --- THE FIX IS HERE ---
  let instructionList = [];
  if (Array.isArray(recipe.instructions)) {
      // Handle Array of Objects (New Backend)
      instructionList = recipe.instructions.map(inst => inst.text || inst);
  } else if (typeof recipe.instructions === 'string') {
      // Handle String (Old/Simple Backend)
      instructionList = recipe.instructions.split('\n').filter(step => step.trim() !== "");
  }
  // -----------------------

  return (
    <div className="min-h-screen bg-orange-50 p-6 flex justify-center">
      <div className="max-w-3xl w-full bg-white p-8 rounded-xl shadow-lg border border-gray-100">
        
        <button onClick={() => navigate('/')} className="flex items-center gap-2 text-gray-500 hover:text-orange-600 mb-6">
            <ArrowLeft size={18} /> Back to Recipes
        </button>

        <div className="flex justify-between items-start mb-4">
            <h1 className="text-3xl font-bold text-gray-900">{recipe.name}</h1>
            <span className="bg-orange-100 text-orange-700 px-3 py-1 rounded-full text-sm font-semibold">
                {recipe.category || "Dinner"}
            </span>
        </div>

        <p className="text-gray-600 mb-6 text-lg italic border-l-4 border-orange-500 pl-4 bg-gray-50 p-2">
            {recipe.description || "No description provided."}
        </p>

        <div className="flex gap-6 text-sm text-gray-600 mb-8 bg-orange-50 p-4 rounded-lg">
          <div className="flex items-center gap-2">
             <Clock size={18} /> 
             <span className="font-semibold">Prep: {recipe.prep_time}m</span>
          </div>
          <div className="flex items-center gap-2">
             <Clock size={18} /> 
             <span className="font-semibold">Cook: {recipe.cook_time}m</span>
          </div>
          <div className="flex items-center gap-2">
             <Users size={18} /> 
             <span className="font-semibold">Servings: {recipe.servings}</span>
          </div>
        </div>

        {/* Ingredients */}
        <h2 className="text-xl font-bold mb-4 text-gray-800 border-b pb-2">Ingredients</h2>
        <ul className="list-disc ml-6 mb-8 space-y-2 text-gray-700">
          {recipe.ingredients && recipe.ingredients.length > 0 ? (
            recipe.ingredients.map((ing, index) => (
              <li key={index}>
                <span className="font-semibold">
                    {ing.name || (ing.ingredient && ing.ingredient.name) || "Item"}
                </span> 
                <span className="text-gray-500"> - {ing.quantity || ing.amount} {ing.unit}</span>
              </li>
            ))
          ) : (
            <p className="text-gray-400 italic">No ingredients listed.</p>
          )}
        </ul>

        {/* Instructions */}
        <h2 className="text-xl font-bold mb-4 text-gray-800 border-b pb-2">Instructions</h2>
        {instructionList.length > 0 ? (
          <div className="space-y-4">
             {instructionList.map((step, index) => (
                 <div key={index} className="flex gap-4">
                     <span className="flex-shrink-0 w-8 h-8 bg-orange-100 text-orange-600 rounded-full flex items-center justify-center font-bold text-sm">
                        {index + 1}
                     </span>
                     <p className="text-gray-700 mt-1">{step}</p>
                 </div>
             ))}
          </div>
        ) : (
          <p className="text-gray-400 italic">No instructions added.</p>
        )}

        <div className="flex justify-end gap-4 mt-10 pt-6 border-t">
          <button
            onClick={() => navigate(`/add-recipe?edit=${recipe.id || recipe.ID}`)}
            className="flex items-center gap-2 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
          >
            <Edit size={18}/> Edit
          </button>
          
          <button
            onClick={handleDelete}
            className="flex items-center gap-2 bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition"
          >
            <Trash2 size={18}/> Delete
          </button>
        </div>
      </div>
    </div>
  );
};

export default RecipeDetails;