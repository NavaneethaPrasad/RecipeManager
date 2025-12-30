import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar';
import api from '../api/axios';
import { Scale, ChefHat, ArrowRight, Calculator } from 'lucide-react';
import toast from 'react-hot-toast';
import EmptyState from '../components/EmptyState'; // Added Import

const ScaleRecipes = () => {
    const [recipes, setRecipes] = useState([]);
    const [selectedRecipeId, setSelectedRecipeId] = useState('');
    const [originalRecipe, setOriginalRecipe] = useState(null); 
    const [targetServings, setTargetServings] = useState(4);
    const [loading, setLoading] = useState(true); // Added loading state
    
    useEffect(() => {
        setLoading(true);
        api.get('/recipes')
            .then(res => {
                setRecipes(res.data);
                setLoading(false);
            })
            .catch(err => {
                console.error("Failed to load recipes");
                setLoading(false);
            });
    }, []);

    useEffect(() => {
        if (!selectedRecipeId) {
            setOriginalRecipe(null);
            return;
        }

        const fetchDetails = async () => {
            try {
                const res = await api.get(`/recipes/${selectedRecipeId}`);
                setOriginalRecipe(res.data);
                setTargetServings(res.data.servings); 
            } catch (err) {
                toast.error("Failed to load ingredients");
            }
        };
        fetchDetails();
    }, [selectedRecipeId]);

    // --- ADDED: Check for Empty State ---
    if (!loading && (!recipes || recipes.length === 0)) {
        return (
            <div className="min-h-screen bg-gray-50">
                <Navbar />
                <div className="container mx-auto p-6 max-w-4xl">
                    <EmptyState 
                        title="Nothing to Scale" 
                        message="Please add a recipe first. Once you have a recipe in your collection, you can use this tool to instantly adjust ingredient quantities for any number of servings." 
                    />
                </div>
            </div>
        );
    }
    // ---------------------------------

    const calculateScaledIngredients = () => {
        if (!originalRecipe || !originalRecipe.ingredients) return [];

        const baseServings = originalRecipe.servings || 1;
        const desiredServings = parseFloat(targetServings) || 0;

        if (desiredServings <= 0) return originalRecipe.ingredients;

        const ratio = desiredServings / baseServings;

        return originalRecipe.ingredients.map(ing => ({
            ...ing,
            scaledQuantity: (ing.quantity * ratio).toFixed(2)
        }));
    };

    const scaledIngredients = calculateScaledIngredients();

    return (
        <div className="min-h-screen bg-gray-50">
            <Navbar />
            
            <div className="container mx-auto p-6 max-w-5xl">
                <div className="flex items-center gap-3 mb-8">
                    <div className="bg-orange-100 p-3 rounded-full">
                        <Scale className="text-orange-600" size={32} />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">Scale Ingredients</h1>
                        <p className="text-gray-500 text-sm">Adjust ingredient quantities instantly</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    
                    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200 h-fit">
                        <h3 className="font-bold text-gray-700 mb-4 flex items-center gap-2">
                            <Calculator size={18} /> Configuration
                        </h3>
                        
                        <div className="space-y-4">
                            <div>
                                <label className="block text-xs font-bold text-gray-500 uppercase mb-1">Recipe</label>
                                <select 
                                    className="w-full p-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 outline-none bg-white"
                                    value={selectedRecipeId}
                                    onChange={(e) => setSelectedRecipeId(e.target.value)}
                                >
                                    <option value="">-- Select Recipe --</option>
                                    {recipes.map(r => (
                                        <option key={r.id} value={r.id}>{r.name}</option>
                                    ))}
                                </select>
                            </div>

                            <div>
                                <label className="block text-xs font-bold text-gray-500 uppercase mb-1">Desired Servings</label>
                                <input 
                                    type="number" 
                                    min="1"
                                    className="w-full p-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 outline-none"
                                    value={targetServings}
                                    onChange={(e) => setTargetServings(e.target.value)}
                                    disabled={!originalRecipe}
                                />
                            </div>
                        </div>
                    </div>

                    <div className="md:col-span-2">
                        {!originalRecipe ? (
                            <div className="bg-white p-12 rounded-xl shadow-sm border border-gray-200 text-center border-dashed">
                                <ChefHat className="mx-auto text-gray-300 mb-3" size={48} />
                                <p className="text-gray-500">Select a recipe to start scaling.</p>
                            </div>
                        ) : (
                            <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                                <div className="p-6 bg-orange-50 border-b border-orange-100 flex justify-between items-center">
                                    <div>
                                        <h2 className="text-xl font-bold text-gray-800">{originalRecipe.name}</h2>
                                        <div className="flex items-center gap-2 text-sm text-gray-600 mt-1">
                                            <span className="bg-white px-2 py-0.5 rounded border">Original: {originalRecipe.servings}</span>
                                            <ArrowRight size={14} />
                                            <span className="bg-orange-600 text-white px-2 py-0.5 rounded font-bold">New: {targetServings}</span>
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        <span className="block text-xs text-gray-500 uppercase font-bold">Ratio</span>
                                        <span className="text-2xl font-bold text-orange-600">
                                            {(targetServings / originalRecipe.servings).toFixed(2)}x
                                        </span>
                                    </div>
                                </div>

                                <div className="p-0">
                                    <table className="w-full text-left border-collapse">
                                        <thead className="bg-gray-50 text-gray-500 text-xs uppercase">
                                            <tr>
                                                <th className="p-4 font-semibold">Ingredient</th>
                                                <th className="p-4 font-semibold text-right">Original</th>
                                                <th className="p-4 font-semibold text-right text-orange-600">Scaled</th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-gray-100">
                                            {scaledIngredients.map((ing, i) => (
                                                <tr key={i} className="hover:bg-gray-50">
                                                    <td className="p-4 font-medium text-gray-800">
                                                        {ing.name || (ing.ingredient && ing.ingredient.name)}
                                                    </td>
                                                    <td className="p-4 text-right text-gray-400">
                                                        {ing.quantity} {ing.unit}
                                                    </td>
                                                    <td className="p-4 text-right font-bold text-orange-600 bg-orange-50/30">
                                                        {ing.scaledQuantity} {ing.unit}
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        )}
                    </div>

                </div>
            </div>
        </div>
    );
};

export default ScaleRecipes;