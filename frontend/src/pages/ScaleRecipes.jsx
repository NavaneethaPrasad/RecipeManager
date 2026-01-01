import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar';
import api from '../api/axios';
import { Scale, ChefHat, ArrowRight, Calculator } from 'lucide-react';
import toast from 'react-hot-toast';
import EmptyState from '../components/EmptyState';

const ScaleRecipes = () => {
    const [recipes, setRecipes] = useState([]);
    const [selectedRecipeId, setSelectedRecipeId] = useState('');
    const [originalRecipe, setOriginalRecipe] = useState(null); 
    const [targetServings, setTargetServings] = useState(4);
    const [loading, setLoading] = useState(true);
    
    useEffect(() => {
        setLoading(true);
        api.get('/recipes').then(res => { setRecipes(res.data || []); setLoading(false); })
            .catch(() => setLoading(false));
    }, []);

    useEffect(() => {
        if (!selectedRecipeId) { setOriginalRecipe(null); return; }
        api.get(`/recipes/${selectedRecipeId}`).then(res => {
            setOriginalRecipe(res.data);
            setTargetServings(res.data.servings);
        });
    }, [selectedRecipeId]);

    const calculateScaledIngredients = () => {
        if (!originalRecipe || !originalRecipe.ingredients) return [];
        const base = originalRecipe.servings || 1;
        const desired = parseFloat(targetServings) || 0;
        if (desired <= 0) return originalRecipe.ingredients;
        const ratio = desired / base;
        return originalRecipe.ingredients.map(ing => ({ 
            ...ing, 
            scaledQuantity: (ing.quantity * ratio).toFixed(1) 
        }));
    };

    if (!loading && recipes.length === 0) {
        return (
            <div className="min-h-screen bg-slate-100">
                <Navbar />
                <div className="container mx-auto p-4 md:p-6 max-w-4xl">
                    <EmptyState title="Nothing to Scale" message="Add a recipe first." />
                </div>
            </div>
        );
    }

    const scaledIngredients = calculateScaledIngredients();

    return (
        <div className="min-h-screen bg-slate-100 pb-10">
            <Navbar />
            <div className="container mx-auto p-4 md:p-6 max-w-6xl">
                
                <div className="mb-6 md:mb-8 bg-white p-4 md:p-6 rounded-2xl shadow-sm border border-gray-100">
                    <h1 className="text-xl md:text-2xl font-black flex items-center gap-3 text-slate-800">
                        <Scale className="text-orange-600" size={28} /> Scale Ingredients
                    </h1>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 md:gap-8">
                    
                    <div className="bg-white p-5 md:p-6 rounded-2xl shadow-sm border border-gray-100 h-fit">
                        <p className="text-[10px] md:text-xs font-black text-slate-500 uppercase mb-4 ml-1 tracking-widest">
                            Configuration
                        </p>
                        <div className="space-y-5 md:space-y-6">
                            <div>
                                <label className="text-xs font-black text-slate-400 uppercase ml-1">Choose Recipe</label>
                                <select 
                                    value={selectedRecipeId} 
                                    onChange={(e) => setSelectedRecipeId(e.target.value)} 
                                    className="w-full p-3 bg-gray-50 border border-gray-100 rounded-xl font-bold text-base md:text-lg text-slate-700 mt-1 outline-none focus:ring-2 focus:ring-orange-500 transition-all"
                                >
                                    <option value="">-- Select --</option>
                                    {recipes.map(r => <option key={r.id} value={r.id}>{r.name}</option>)}
                                </select>
                            </div>
                            <div>
                                <label className="text-xs font-black text-slate-400 uppercase ml-1">Desired Servings</label>
                                <input 
                                    type="number" 
                                    value={targetServings} 
                                    onChange={(e) => setTargetServings(e.target.value)} 
                                    className="w-full p-3 bg-gray-50 border border-gray-100 rounded-xl font-bold text-base md:text-lg text-slate-700 mt-1 outline-none focus:ring-2 focus:ring-orange-500 transition-all" 
                                />
                            </div>
                        </div>
                    </div>

                    <div className="lg:col-span-2 bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden min-h-[400px] md:min-h-[500px]">
                        {!originalRecipe ? (
                            <div className="p-16 md:p-20 text-center text-slate-400 font-bold text-base md:text-lg">
                                Select a recipe to view results.
                            </div>
                        ) : (
                            <div className="overflow-x-auto"> 
                                <table className="w-full border-collapse">
                                    <thead className="bg-gray-50 border-b border-gray-100">
                                        <tr className="text-left text-[10px] md:text-xs font-black text-slate-400 uppercase">
                                            <th className="p-4 md:p-5 tracking-tight">Ingredient</th>
                                            <th className="p-4 md:p-5 text-right tracking-tight">Original</th>
                                            <th className="p-4 md:p-5 text-right text-orange-600 tracking-tight">Scaled</th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-50">
                                        {scaledIngredients.map((ing, i) => (
                                            <tr key={i} className="hover:bg-gray-50/50 transition-colors">
                                                <td className="p-4 md:p-5 font-bold text-base md:text-lg text-slate-700">
                                                    {ing.name || ing.ingredient?.name}
                                                </td>
                                                <td className="p-4 md:p-5 text-right text-slate-400 font-bold text-sm">
                                                    {ing.quantity} <span className="text-[10px] uppercase font-black">{ing.unit}</span>
                                                </td>
                                                <td className="p-4 md:p-5 text-right font-black text-lg md:text-xl text-orange-600 bg-orange-50/20">
                                                    {ing.scaledQuantity} <span className="text-[10px] uppercase">{ing.unit}</span>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ScaleRecipes;