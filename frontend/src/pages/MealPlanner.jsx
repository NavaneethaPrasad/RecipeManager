import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar';
import AddMealModal from '../components/AddMealModal'; 
import api from '../api/axios';
import { ChevronLeft, ChevronRight, Plus, Trash2, Calendar } from 'lucide-react';
import toast from 'react-hot-toast';
import EmptyState from '../components/EmptyState';

const MealPlanner = () => {
    const [currentDate, setCurrentDate] = useState(new Date());
    const [weekDates, setWeekDates] = useState([]);
    const [plans, setPlans] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [modalData, setModalData] = useState({ date: '', type: '' });
    const [hasRecipes, setHasRecipes] = useState(true);

    useEffect(() => {
        const startOfWeek = new Date(currentDate);
        const day = startOfWeek.getDay(); 
        const diff = startOfWeek.getDate() - day + (day === 0 ? -6 : 1); 
        startOfWeek.setDate(diff);

        const days = [];
        for (let i = 0; i < 7; i++) {
            const d = new Date(startOfWeek);
            d.setDate(startOfWeek.getDate() + i);
            days.push(d);
        }
        setWeekDates(days);

        api.get('/recipes').then(res => {
            setHasRecipes((res.data || []).length > 0);
        });

        if(days.length > 0) fetchPlans(days[0], days[6]);
    }, [currentDate]);

    const fetchPlans = async (start, end) => {
        try {
            const s = start.toISOString().split('T')[0];
            const e = end.toISOString().split('T')[0];
            const res = await api.get(`/meal-plans?start_date=${s}&end_date=${e}`);
            setPlans(res.data || []);
        } catch (err) { console.error("Failed to fetch plans"); }
    };

    const changeWeek = (offset) => {
        const d = new Date(currentDate); 
        d.setDate(d.getDate() + (offset * 7)); 
        setCurrentDate(d);
    };

    const openAddModal = (date, type) => {
        setModalData({ date: date.toISOString().split('T')[0], type });
        setIsModalOpen(true);
    };

    const handleDelete = async (id) => {
        if(!confirm("Remove meal?")) return;
        try {
            await api.delete(`/meal-plans/${id}`);
            toast.success("Meal removed");
            fetchPlans(weekDates[0], weekDates[6]);
        } catch (err) { toast.error("Delete failed"); }
    };

    const getPlan = (date, type) => {
        const dateStr = date.toISOString().split('T')[0];
        return (plans || []).find(p => p.date && p.date.startsWith(dateStr) && p.meal_type === type);
    };

    const MEAL_TYPES = ['breakfast', 'lunch', 'dinner', 'snack'];

    if (!hasRecipes) {
        return (
            <div className="min-h-screen bg-slate-100">
                <Navbar />
                <div className="container mx-auto p-4 md:p-6 max-w-4xl">
                    <EmptyState title="Add Recipes First" message="You need recipes in your collection to start planning your weekly meals." />
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-slate-100 pb-10">
            <Navbar />
            <div className="container mx-auto p-4 md:p-6 max-w-7xl">
                
                <div className="flex flex-col lg:flex-row justify-between items-center gap-4 mb-8 bg-white p-4 md:p-6 rounded-2xl shadow-sm border border-gray-100">
                    <h1 className="text-xl md:text-2xl font-black flex items-center gap-3 text-slate-800">
                        <Calendar className="text-orange-600" size={28} /> Weekly Planner
                    </h1>
                    
                    <div className="flex gap-2 md:gap-4 items-center bg-gray-50 p-2 rounded-xl border border-gray-100 w-full lg:w-auto justify-between">
                        <button onClick={() => changeWeek(-1)} className="p-2 bg-white rounded-lg shadow-xs hover:bg-gray-50"><ChevronLeft size={20}/></button>
                        <span className="font-bold text-slate-700 text-xs md:text-sm text-center px-2">
                            {weekDates[0]?.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })} - {weekDates[6]?.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })}
                        </span>
                        <button onClick={() => changeWeek(1)} className="p-2 bg-white rounded-lg shadow-xs hover:bg-gray-50"><ChevronRight size={20}/></button>
                    </div>
                </div>

                <div className="overflow-x-auto pb-6 -mx-4 px-4 md:mx-0 md:px-0">
                    <div className="grid grid-cols-7 gap-3 md:gap-4 min-w-[1000px] md:min-w-full min-h-[550px]">
                        {weekDates.map((date, i) => (
                            <div key={i} className="flex flex-col h-full bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                                <div className={`p-3 md:p-4 text-center border-b-2 ${date.toDateString() === new Date().toDateString() ? 'bg-orange-50 border-orange-500' : 'bg-gray-50/50 border-gray-100'}`}>
                                    <p className="font-black text-[10px] md:text-sm uppercase text-slate-400 mb-1">{date.toLocaleDateString(undefined, {weekday:'short'})}</p>
                                    <p className="text-lg md:text-2xl font-black text-slate-800">{date.getDate()}</p>
                                </div>
                                <div className="p-3 md:p-4 flex flex-col gap-4 md:gap-6 flex-1">
                                    {MEAL_TYPES.map(type => {
                                        const plan = getPlan(date, type);
                                        return (
                                            <div key={type}>
                                                <p className="text-[10px] md:text-sm font-black text-slate-500 uppercase mb-2 md:mb-3 ml-1 tracking-tighter">{type}</p>
                                                <div className="h-20 md:h-24 relative group">
                                                    {plan ? (
                                                        <div className="h-full bg-orange-100/50 p-2 md:p-4 rounded-xl border border-orange-200 flex flex-col justify-center transition-all group-hover:bg-orange-100">
                                                            <p className="font-bold text-slate-800 text-xs md:text-lg leading-tight line-clamp-2">{plan.recipe?.name}</p>
                                                            <p className="text-[9px] md:text-[11px] font-black text-orange-600 uppercase mt-1">
                                                                {plan.target_servings} Servings
                                                            </p>
                                                            <button onClick={() => handleDelete(plan.id || plan.ID)} className="absolute -top-1 -right-1 bg-red-500 text-white p-1 rounded-full opacity-0 group-hover:opacity-100 transition-opacity shadow-lg"><Trash2 size={12}/></button>
                                                        </div>
                                                    ) : (
                                                        <button onClick={() => openAddModal(date, type)} className="w-full h-full border-2 border-dashed border-slate-100 rounded-xl flex items-center justify-center text-slate-300 hover:text-orange-500 hover:bg-orange-50 transition-all"><Plus size={24}/></button>
                                                    )}
                                                </div>
                                            </div>
                                        )
                                    })}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
            {isModalOpen && (
                <AddMealModal 
                    date={modalData.date} mealType={modalData.type} 
                    onClose={() => setIsModalOpen(false)} 
                    onSave={() => fetchPlans(weekDates[0], weekDates[6])} 
                />
            )}
        </div>
    );
};
export default MealPlanner;