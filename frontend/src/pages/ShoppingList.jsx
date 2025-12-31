import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar';
import api from '../api/axios';
import { ShoppingCart, RefreshCw, CheckSquare, Square, Calendar } from 'lucide-react';
import toast from 'react-hot-toast';
import EmptyState from '../components/EmptyState';

const ShoppingList = () => {
    const [shoppingList, setShoppingList] = useState(null);
    const [loading, setLoading] = useState(false);
    const [hasRecipes, setHasRecipes] = useState(true);
    
    const today = new Date();
    const nextWeek = new Date();
    nextWeek.setDate(today.getDate() + 7);

    const [dates, setDates] = useState({
        start: today.toISOString().split('T')[0],
        end: nextWeek.toISOString().split('T')[0]
    });

    useEffect(() => {
        api.get('/recipes').then(res => setHasRecipes((res.data || []).length > 0));
        const savedId = localStorage.getItem('last_shopping_list_id');
        if (savedId) fetchList(savedId);
    }, []);

    const fetchList = async (id) => {
        setLoading(true);
        try {
            const res = await api.get(`/shopping-lists/${id}`);
            setShoppingList(res.data);
        } catch (err) { localStorage.removeItem('last_shopping_list_id'); }
        finally { setLoading(false); }
    };

    const handleGenerate = async () => {
        setLoading(true);
        try {
            const res = await api.post('/shopping-lists/generate', { start_date: dates.start, end_date: dates.end });
            setShoppingList(res.data);
            if (res.data?.id) localStorage.setItem('last_shopping_list_id', res.data.id);
            toast.success("Shopping List Generated!");
        } catch (err) { toast.error("Check your meal plan first!"); }
        finally { setLoading(false); }
    };

    const handleToggle = async (itemId) => {
        const updatedItems = (shoppingList?.items || []).map(item => item.id === itemId ? { ...item, checked: !item.checked } : item);
        setShoppingList({ ...shoppingList, items: updatedItems });
        try { await api.patch(`/shopping-lists/items/${itemId}/toggle`); }
        catch (err) { if (shoppingList?.id) fetchList(shoppingList.id); }
    };

    if (!hasRecipes) {
        return (
            <div className="min-h-screen bg-slate-100">
                <Navbar />
                <div className="container mx-auto p-6 max-w-4xl"><EmptyState title="No Recipes Found" message="Add recipes to generate shopping lists." /></div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-slate-100 pb-10">
            <Navbar />
            <div className="container mx-auto p-6 max-w-2xl">
                <div className="flex justify-between items-center mb-8 bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
                    <h1 className="text-2xl font-black flex items-center gap-3 text-slate-800">
                        <ShoppingCart className="text-orange-600" size={28} /> Shopping List
                    </h1>
                    <button onClick={handleGenerate} disabled={loading} className="bg-orange-600 text-white px-6 py-2 rounded-xl font-bold hover:bg-orange-700 transition disabled:opacity-50">
                        {loading ? "..." : "Generate"}
                    </button>
                </div>

                <div className="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 mb-6 flex gap-4">
                     <div className="flex-1">
                        <p className="text-sm font-black text-slate-500 uppercase mb-2 ml-1">Start</p>
                        <input type="date" value={dates.start} onChange={e => setDates({...dates, start: e.target.value})} className="w-full p-2 bg-gray-50 border border-gray-100 rounded-xl font-bold text-lg text-slate-700 outline-none" />
                     </div>
                     <div className="flex-1">
                        <p className="text-sm font-black text-slate-500 uppercase mb-2 ml-1">End</p>
                        <input type="date" value={dates.end} onChange={e => setDates({...dates, end: e.target.value})} className="w-full p-2 bg-gray-50 border border-gray-100 rounded-xl font-bold text-lg text-slate-700 outline-none" />
                     </div>
                </div>

                <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden min-h-[400px]">
                    {(!shoppingList || (shoppingList.items || []).length === 0) ? (
                        <div className="p-20 text-center text-slate-400 font-bold text-lg">No items yet.</div>
                    ) : (
                        <div className="divide-y divide-gray-50">
                            {shoppingList.items.map((item) => (
                                <div key={item.id} onClick={() => handleToggle(item.id)} className="p-5 flex items-center justify-between hover:bg-gray-50 cursor-pointer">
                                    <div className="flex items-center gap-4">
                                        <div className={item.checked ? "text-green-500" : "text-slate-300"}>
                                            {item.checked ? <CheckSquare size={26}/> : <Square size={26}/>}
                                        </div>
                                        <p className={`text-lg font-bold ${item.checked ? 'line-through text-slate-300' : 'text-slate-700'}`}>{item.name}</p>
                                    </div>
                                    <p className="text-xl font-black text-orange-600">
                                        {parseFloat(item.quantity).toFixed(1)} <span className="text-xs text-slate-400 uppercase">{item.unit}</span>
                                    </p>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};
export default ShoppingList;