import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import api from "../api/axios";
import toast from "react-hot-toast";
import Navbar from "../components/Navbar";
import { Edit, Trash2, Clock, Users, ArrowLeft } from "lucide-react";

const RecipeDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [recipe, setRecipe] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get(`/recipes/${id}`)
      .then(res => { setRecipe(res.data); setLoading(false); })
      .catch(() => { toast.error("Recipe not found"); navigate("/recipes"); });
  }, [id, navigate]);

  const handleDelete = async () => {
    if (!window.confirm("Delete this recipe permanently?")) return;
    try {
      await api.delete(`/recipes/${id}`);
      toast.success("Recipe deleted");
      navigate("/recipes");
    } catch { toast.error("Delete failed"); }
  };

  if (loading || !recipe) return <div className="flex h-screen items-center justify-center font-black text-orange-600 bg-slate-100">Loading...</div>;

  const instructionList = Array.isArray(recipe.instructions) ? recipe.instructions : [];

  return (
    <div className="min-h-screen bg-slate-100 pb-10">
      <Navbar />
      <div className="container mx-auto p-6 max-w-4xl">
        
        {/* Header Container */}
        <div className="flex justify-between items-center mb-8 bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
            <div className="flex items-center gap-4">
                <Link to="/recipes" className="p-2 hover:bg-gray-100 rounded-full transition text-slate-400">
                    <ArrowLeft size={24} />
                </Link>
                <h1 className="text-2xl font-black text-slate-800">{recipe.name}</h1>
            </div>
            <div className="flex gap-2">
                <button onClick={() => navigate(`/add-recipe?edit=${recipe.id || recipe.ID}`)} className="p-2 bg-slate-50 text-slate-600 rounded-xl hover:bg-slate-100 transition border border-slate-100"><Edit size={20}/></button>
                <button onClick={handleDelete} className="p-2 bg-red-50 text-red-600 rounded-xl hover:bg-red-100 transition border border-red-100"><Trash2 size={20}/></button>
            </div>
        </div>

        {/* Content Card */}
        <div className="bg-white p-8 rounded-3xl shadow-sm border border-gray-100">
            <div className="grid grid-cols-3 gap-4 mb-8">
                <div className="bg-slate-50 p-4 rounded-2xl border border-slate-100 text-center">
                    <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Preparation</p>
                    <p className="text-xl font-black text-slate-700">{recipe.prep_time} <span className="text-sm">min</span></p>
                </div>
                <div className="bg-slate-50 p-4 rounded-2xl border border-slate-100 text-center">
                    <p className="text-[10px] font-black text-slate-400 uppercase mb-1">Cooking</p>
                    <p className="text-xl font-black text-slate-700">{recipe.cook_time} <span className="text-sm">min</span></p>
                </div>
                <div className="bg-orange-50 p-4 rounded-2xl border border-orange-100 text-center">
                    <p className="text-[10px] font-black text-orange-400 uppercase mb-1">Servings</p>
                    <p className="text-xl font-black text-orange-600">{recipe.servings} <span className="text-sm">ppl</span></p>
                </div>
            </div>

            <p className="text-xl text-slate-600 italic mb-10 leading-relaxed border-l-4 border-orange-500 pl-6 py-2 bg-orange-50/30 rounded-r-xl">
                {recipe.description || "No description provided."}
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
                <div>
                    <h2 className="text-xs font-black text-slate-400 uppercase tracking-widest mb-6">Ingredients</h2>
                    <ul className="space-y-4">
                        {(recipe.ingredients || []).map((ing, index) => (
                            <li key={index} className="flex justify-between items-center bg-gray-50 p-3 rounded-xl">
                                <span className="font-bold text-lg text-slate-700">{ing.name || ing.ingredient?.name}</span>
                                <span className="font-black text-orange-600">{ing.quantity || ing.amount} <span className="text-[10px] uppercase text-slate-400">{ing.unit}</span></span>
                            </li>
                        ))}
                    </ul>
                </div>

                <div>
                    <h2 className="text-xs font-black text-slate-400 uppercase tracking-widest mb-6">Instructions</h2>
                    <div className="space-y-6">
                        {instructionList.map((step, index) => (
                            <div key={index} className="flex gap-4">
                                <span className="flex-shrink-0 w-6 h-6 bg-slate-800 text-white rounded-md flex items-center justify-center font-black text-xs mt-1">
                                    {index + 1}
                                </span>
                                <p className="text-lg font-bold text-slate-600 leading-snug">{step.text || step}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
      </div>
    </div>
  );
};
export default RecipeDetails;