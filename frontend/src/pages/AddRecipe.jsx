import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import api from '../api/axios';
import toast from 'react-hot-toast';
import Navbar from '../components/Navbar';
import { Plus, Minus, Save, ArrowLeft } from 'lucide-react';

const AddRecipe = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const editId = searchParams.get('edit');
    
    const [form, setForm] = useState({
        name: '', description: '', prep_time: 0, cook_time: 0, servings: 4, category: '-Select-', instructions: ''
    });
    const [ingredients, setIngredients] = useState([{ name: '', amount: '', unit: '' }]);
    const [originalData, setOriginalData] = useState(null);

    useEffect(() => {
        if (editId) {
            api.get(`/recipes/${editId}`).then(res => {
                const r = res.data;
                const formattedData = {
                    form: {
                        name: r.name,
                        description: r.description,
                        prep_time: r.prep_time,
                        cook_time: r.cook_time,
                        servings: r.servings,
                        category: r.category,
                        instructions: Array.isArray(r.instructions) ? r.instructions.map(i => i.text).join('\n') : r.instructions
                    },
                    ingredients: (r.ingredients || []).map(ing => ({
                        name: ing.name || ing.ingredient?.name || "",
                        amount: ing.quantity || ing.amount,
                        unit: ing.unit
                    }))
                };
                setForm(formattedData.form);
                setIngredients(formattedData.ingredients);
                setOriginalData(JSON.stringify(formattedData));
            });
        }
    }, [editId]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const payload = {
            ...form,
            prep_time: parseInt(form.prep_time),
            cook_time: parseInt(form.cook_time),
            servings: parseInt(form.servings),
            instructions: form.instructions.split('\n').filter(line => line.trim() !== ""),
            ingredients: ingredients.map(ing => ({ name: ing.name, amount: parseFloat(ing.amount), unit: ing.unit }))
        };

        if (editId && originalData === JSON.stringify({form, ingredients})) {
            toast("No changes detected.", { icon: 'ℹ️' });
            navigate(`/recipes/${editId}`);
            return;
        }

        try {
            if (editId) await api.put(`/recipes/${editId}`, payload);
            else await api.post('/recipes', payload);
            toast.success(editId ? "Updated successfully!" : "Created successfully!");
            navigate('/recipes');
        } catch (err) { toast.error("Failed to save"); }
    };

    const inputClass = "w-full p-3 bg-slate-50 border border-gray-100 rounded-2xl font-bold text-lg text-slate-700 outline-none focus:ring-2 focus:ring-orange-500";
    const labelClass = "text-sm font-black text-slate-500 uppercase mb-2 ml-1 block";

    return (
        <div className="min-h-screen bg-slate-100 pb-10">
            <Navbar />
           <form onSubmit={handleSubmit} className="container mx-auto p-6 max-w-4xl">
                <div className="flex justify-between items-center mb-8 bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
                    <div className="flex items-center gap-4">
                        <Link to="/recipes" className="p-2 hover:bg-gray-100 rounded-full transition text-slate-400"><ArrowLeft /></Link>
                        <h1 className="text-2xl font-black text-slate-800">{editId ? "Edit Recipe" : "New Recipe"}</h1>
                    </div>
                    <button type="submit" className="bg-orange-600 text-white px-8 py-3 rounded-2xl font-black hover:bg-orange-700 transition shadow-lg flex items-center gap-2">
                        <Save size={20}/> {editId ? "Update" : "Save"}
                    </button>
                </div>

                <div className="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 space-y-8">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="md:col-span-2">
                            <label className={labelClass}>Recipe Name</label>
                            <input className={inputClass} value={form.name} onChange={e => setForm({...form, name: e.target.value})} required />
                        </div>
                        <div className="md:col-span-2">
                            <label className={labelClass}>Description</label>
                            <input className={inputClass} value={form.description} onChange={e => setForm({...form, description: e.target.value})} placeholder="e.g. A spicy traditional curry" />
                        </div>
                        <div>
                            <label className={labelClass}>Category</label>
                            <select className={inputClass} value={form.category} onChange={e => setForm({...form, category: e.target.value})} required >
                                <option value="-Select-">-Select-</option><option value="General">General</option><option value="Breakfast">Breakfast</option><option value="Lunch">Lunch</option><option value="Dinner">Dinner</option><option value="Snack">Snack</option>
                            </select>
                        </div>
                        <div>
                            <label className={labelClass}>Servings</label>
                            <input type="number" className={inputClass} value={form.servings} onChange={e => setForm({...form, servings: e.target.value})} required />
                        </div>
                        <div>
                            <label className={labelClass}>Prep Time (m)</label>
                            <input type="number" className={inputClass} value={form.prep_time} onChange={e => setForm({...form, prep_time: e.target.value})} required />
                        </div>
                        <div>
                            <label className={labelClass}>Cook Time (m)</label>
                            <input type="number" className={inputClass} value={form.cook_time} onChange={e => setForm({...form, cook_time: e.target.value})} required />
                        </div>
                    </div>

                    <div>
                        <label className={labelClass}>Ingredients</label>
                        {ingredients.map((ing, i) => (
                            <div key={i} className="flex gap-3 mb-3 items-center">
                                <input placeholder="Ingredient Name" className={`${inputClass} w-24`} value={ing.name} onChange={e => {const l=[...ingredients]; l[i].name=e.target.value; setIngredients(l);}} required />
                                <input placeholder="Qty" type="number" className={`${inputClass} w-24`} value={ing.amount} onChange={e => {const l=[...ingredients]; l[i].amount=e.target.value; setIngredients(l);}} required />
                                <input placeholder="Unit" className={`${inputClass} w-24`} value={ing.unit} onChange={e => {const l=[...ingredients]; l[i].unit=e.target.value; setIngredients(l);}} required />
                                <button type="button" onClick={() => {const l=[...ingredients]; l.splice(i,1); setIngredients(l);}} className="text-red-400 p-2"><Minus/></button>
                            </div>
                        ))}
                        <button type="button" onClick={() => setIngredients([...ingredients, {name:'',amount:'',unit:''}])} className="text-orange-600 font-black text-sm uppercase mt-2">+ Add Ingredient</button>
                    </div>

                    <div>
                        <label className={labelClass}>Instructions</label>
                        <textarea className={`${inputClass} h-40`} placeholder="One step per line..." value={form.instructions} onChange={e => setForm({...form, instructions: e.target.value})} required />
                    </div>
                </div>
             </form> 
        </div>
    );
};
export default AddRecipe;