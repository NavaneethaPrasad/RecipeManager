import { useState, useEffect } from 'react';
import api from '../api/axios';
import { X, Calendar, Users } from 'lucide-react';
import toast from 'react-hot-toast';

const AddMealModal = ({ date, mealType, onClose, onSave }) => {
    const [recipes, setRecipes] = useState([]);
    const [selectedId, setSelectedId] = useState('');
    const [servings, setServings] = useState(2);

    useEffect(() => {
        api.get('/recipes').then(res => setRecipes(res.data || [])).catch(console.error);
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await api.post('/meal-plans', { date, meal_type: mealType, recipe_id: parseInt(selectedId), target_servings: parseInt(servings) });
            toast.success("Planned!"); onSave(); onClose();
        } catch { toast.error("Select a recipe"); }
    };

    return (
        <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in duration-200">
                <div className="bg-orange-600 p-6 text-white flex justify-between items-center">
                    <div>
                        <h2 className="text-2xl font-black uppercase tracking-tighter">Add {mealType}</h2>
                        <p className="text-orange-100 font-bold text-sm">{new Date(date).toDateString()}</p>
                    </div>
                    <button onClick={onClose} className="p-2 hover:bg-orange-700 rounded-xl transition"><X/></button>
                </div>
                
                <form onSubmit={handleSubmit} className="p-8 space-y-6">
                    <div>
                        <label className="text-xs font-black text-slate-400 uppercase mb-2 block ml-1">Pick a Recipe</label>
                        <select className="w-full p-4 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-lg outline-none focus:ring-2 focus:ring-orange-500"
                            onChange={e => setSelectedId(e.target.value)} required>
                            <option value="">-- Choose Recipe --</option>
                            {recipes.map(r => <option key={r.id} value={r.id}>{r.name}</option>)}
                        </select>
                    </div>
                    <div>
                        <label className="text-xs font-black text-slate-400 uppercase mb-2 block ml-1 text-lg">Servings</label>
                        <input type="number" className="w-full p-4 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-lg outline-none" 
                            value={servings} onChange={e => setServings(e.target.value)} min="1" />
                    </div>
                    <button className="w-full bg-slate-800 text-white py-4 rounded-2xl font-black text-lg hover:bg-slate-900 transition-all shadow-lg">Confirm Plan</button>
                </form>
            </div>
        </div>
    );
};
export default AddMealModal;