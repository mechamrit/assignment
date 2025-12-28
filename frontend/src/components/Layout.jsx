import React from 'react';
import { useAuth } from '../context/AuthContext';
import { LogOut, User, ClipboardList } from 'lucide-react';

const Layout = ({ children }) => {
    const { user, logout } = useAuth();

    return (
        <div className="min-h-screen bg-gray-950 text-gray-100 flex flex-col">
            <nav className="bg-gray-900 border-b border-gray-800 px-6 py-4 flex items-center justify-between shadow-md">
                <div className="flex items-center space-x-2">
                    <ClipboardList className="text-blue-500 w-8 h-8" />
                    <span className="text-xl font-bold tracking-tight">QC System</span>
                </div>

                {user && (
                    <div className="flex items-center space-x-6">
                        <div className="flex items-center space-x-2 text-sm">
                            <User className="w-4 h-4 text-gray-400" />
                            <span className="font-medium">{user.role.toUpperCase()}</span>
                        </div>
                        <button
                            onClick={logout}
                            className="flex items-center space-x-2 bg-gray-800 hover:bg-gray-700 px-4 py-2 rounded-lg transition-all text-sm font-medium border border-gray-700"
                        >
                            <LogOut className="w-4 h-4" />
                            <span>Sign Out</span>
                        </button>
                    </div>
                )}
            </nav>

            <main className="flex-1 container mx-auto px-6 py-8">
                {children}
            </main>

            <footer className="bg-gray-900 border-t border-gray-800 py-4 text-center text-xs text-gray-500">
                &copy; 2025 Engineering Quality Control System. Secure Workflow Enabled.
            </footer>
        </div>
    );
};

export default Layout;
