import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authService } from '../services/authService';
import { UserPlus } from 'lucide-react';

const Register = () => {
    const [formData, setFormData] = useState({
        username: '',
        password: '',
        role: 'drafter'
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (loading || success) return;

        setLoading(true);
        setError('');
        setSuccess('');
        try {
            await authService.register(formData);
            setSuccess('Registration successful! Redirecting to login...');
            setTimeout(() => navigate('/login'), 2000);
        } catch (err) {
            setError(err.response?.data?.error || 'Registration failed');
        } finally {
            setLoading(false);
        }
    };

    const roles = [
        { value: 'admin', label: 'Admin' },
        { value: 'drafter', label: 'Drafter' },
        { value: 'shift_lead', label: 'Shift Lead' },
        { value: 'final_qc', label: 'Final QC' }
    ];

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-950 text-white">
            <div className="bg-gray-900 p-8 rounded-xl shadow-2xl w-full max-w-md border border-gray-800">
                <div className="flex items-center justify-center mb-8">
                    <div className="bg-blue-600/20 p-3 rounded-lg mr-3">
                        <UserPlus className="w-8 h-8 text-blue-500" />
                    </div>
                    <h1 className="text-3xl font-extrabold tracking-tight">Create Account</h1>
                </div>

                <form onSubmit={handleSubmit} className="space-y-5">
                    <div>
                        <label className="block text-sm font-semibold text-gray-400 mb-1.5 ml-1">Username</label>
                        <input
                            type="text"
                            name="username"
                            className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all placeholder-gray-500"
                            placeholder="Engineering ID or name"
                            value={formData.username}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-semibold text-gray-400 mb-1.5 ml-1">Password</label>
                        <input
                            type="password"
                            name="password"
                            className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all placeholder-gray-500"
                            placeholder="Minimum 8 characters"
                            value={formData.password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-semibold text-gray-400 mb-1.5 ml-1">System Role</label>
                        <select
                            name="role"
                            className="w-full bg-gray-800 border border-gray-700 rounded-xl px-4 py-3 focus:ring-2 focus:ring-blue-500 outline-none cursor-pointer appearance-none"
                            value={formData.role}
                            onChange={handleChange}
                            required
                        >
                            {roles.map(r => (
                                <option key={r.value} value={r.value} className="bg-gray-900">
                                    {r.label}
                                </option>
                            ))}
                        </select>
                    </div>

                    {error && (
                        <div className="bg-red-900/20 border border-red-500/50 text-red-400 text-sm p-3 rounded-lg flex items-center">
                            <span className="mr-2 italic">⚠️</span> {error}
                        </div>
                    )}

                    {success && (
                        <div className="bg-green-900/20 border border-green-500/50 text-green-400 text-sm p-3 rounded-lg">
                            {success}
                        </div>
                    )}

                    <button
                        type="submit"
                        disabled={loading || success}
                        className={`w-full font-bold py-3.5 rounded-xl transition-all shadow-lg transform ${loading || success
                                ? 'bg-gray-700 cursor-not-allowed text-gray-400'
                                : 'bg-blue-600 hover:bg-blue-700 text-white shadow-blue-900/20 hover:-translate-y-0.5'
                            }`}
                    >
                        {loading ? 'Registering...' : success ? 'Successfully Registered!' : 'Register Member'}
                    </button>
                </form>

                <div className="mt-8 text-center text-gray-400">
                    <p className="text-sm">
                        Already have an account?{' '}
                        <Link to="/login" className="text-blue-400 hover:text-blue-300 font-bold transition-colors">
                            Sign In
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default Register;
