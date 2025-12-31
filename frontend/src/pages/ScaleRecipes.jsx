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
        return originalRecipe.ingredients.map(ing => ({ ...ing, scaledQuantity: (ing.quantity * ratio).toFixed(1) }));
    };

    if (!loading && recipes.length === 0) {
        return <div className="min-h-screen bg-slate-100"><Navbar /><div className="container mx-auto p-6 max-w-4xl"><EmptyState title="Nothing to Scale" message="Add a recipe first." /></div></div>;
    }

    const scaledIngredients = calculateScaledIngredients();

    return (
        <div className="min-h-screen bg-slate-100 pb-10">
            <Navbar />
            <div className="container mx-auto p-6 max-w-6xl">
                <div className="mb-8 bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
                    <h1 className="text-2xl font-black flex items-center gap-3 text-slate-800">
                        <Scale className="text-orange-600" size={28} /> Scale Ingredients
                    </h1>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 h-fit">
                        <p className="text-sm font-black text-slate-500 uppercase mb-4 ml-1">Settings</p>
                        <div className="space-y-6">
                            <div>
                                <label className="text-sm font-black text-slate-400 uppercase ml-1">Recipe</label>
                                <select value={selectedRecipeId} onChange={(e) => setSelectedRecipeId(e.target.value)} className="w-full p-3 bg-gray-50 border border-gray-100 rounded-xl font-bold text-lg text-slate-700 mt-1 outline-none">
                                    <option value="">-- Choose --</option>
                                    {recipes.map(r => <option key={r.id} value={r.id}>{r.name}</option>)}
                                </select>
                            </div>
                            <div>
                                <label className="text-sm font-black text-slate-400 uppercase ml-1">Servings</label>
                                <input type="number" value={targetServings} onChange={(e) => setTargetServings(e.target.value)} className="w-full p-3 bg-gray-50 border border-gray-100 rounded-xl font-bold text-lg text-slate-700 mt-1 outline-none" />
                            </div>
                        </div>
                    </div>

                    <div className="md:col-span-2 bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden min-h-[500px]">
                        {!originalRecipe ? (
                            <div className="p-20 text-center text-slate-400 font-bold text-lg">Select a recipe to view results.</div>
                        ) : (
                            <table className="w-full">
                                <thead className="bg-gray-50 border-b border-gray-100">
                                    <tr className="text-left text-sm font-black text-slate-400 uppercase">
                                        <th className="p-5 tracking-tight">Ingredient</th>
                                        <th className="p-5 text-right tracking-tight">Original</th>
                                        <th className="p-5 text-right text-orange-600 tracking-tight">Scaled</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-gray-50">
                                    {scaledIngredients.map((ing, i) => (
                                        <tr key={i} className="hover:bg-gray-50/50 transition">
                                            <td className="p-5 font-bold text-lg text-slate-700">{ing.name || ing.ingredient?.name}</td>
                                            <td className="p-5 text-right text-slate-400 font-bold">{ing.quantity} {ing.unit}</td>
                                            <td className="p-5 text-right font-black text-xl text-orange-600 bg-orange-50/20">{ing.scaledQuantity} {ing.unit}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};
export default ScaleRecipes;