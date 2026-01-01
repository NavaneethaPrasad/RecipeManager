import { Link } from 'react-router-dom';
import { ChefHat, Plus } from 'lucide-react';

const EmptyState = ({ title, message, showAddButton = true }) => {
    return (
        <div className="flex flex-col items-center justify-center py-12 px-6 md:py-24 bg-white rounded-[2.5rem] border-2 border-dashed border-slate-200 mt-4 md:mt-8 shadow-sm text-center">
            <div className="bg-orange-100 p-4 md:p-6 rounded-3xl mb-6 animate-bounce duration-[3000ms]">
                <ChefHat size={40} className="text-orange-600 md:w-12 md:h-12" />
            </div>
            <h2 className="text-xl md:text-2xl font-black text-slate-800 mb-3 tracking-tight">
                {title}
            </h2>
            <p className="text-slate-500 font-medium text-base md:text-lg max-w-xs md:max-w-md mb-8 leading-relaxed">
                {message}
            </p>

            {showAddButton && (
                <Link 
                    to="/add-recipe" 
                    className="w-full sm:w-auto bg-orange-600 text-white px-8 py-4 rounded-2xl font-black text-lg hover:bg-orange-700 transition-all flex items-center justify-center gap-2 shadow-lg shadow-orange-100 active:scale-95"
                >
                    <Plus size={22} strokeWidth={3} /> 
                    <span>Add Your First Recipe</span>
                </Link>
            )}
        </div>
    );
};

export default EmptyState;