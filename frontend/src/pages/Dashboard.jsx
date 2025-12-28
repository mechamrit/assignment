import React, { useEffect, useState, useRef } from 'react';
import { drawingService } from '../services/drawingService';
import { useAuth } from '../context/AuthContext';
import { Clock, UserPlus } from 'lucide-react';
import TaskCard from '../components/TaskCard';
import OverviewTable from '../components/OverviewTable';
import ProjectSelector from '../components/ProjectSelector';

const Dashboard = () => {
    const [drawings, setDrawings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [currentProjectID, setCurrentProjectID] = useState(1); // Default to Project 1
    const { user } = useAuth();
    const eventSourceRef = useRef(null);

    const fetchDrawings = async (projectID) => {
        setLoading(true);
        try {
            const data = await drawingService.getAll(projectID);
            setDrawings(data);
            setError('');
        } catch (err) {
            setError('Failed to load drawings');
        } finally {
            setLoading(false);
        }
    };

    // Initial fetch and SSE setup
    useEffect(() => {
        fetchDrawings(currentProjectID);

        // Close existing connection if any
        if (eventSourceRef.current) {
            eventSourceRef.current.close();
        }

        const token = localStorage.getItem('token');
        const eventSource = new EventSource(`/api/v1/events?project_id=${currentProjectID}&token=${token}`);
        eventSourceRef.current = eventSource;

        eventSource.onmessage = (event) => {
            const data = JSON.parse(event.data);
            // Optimization: We could merge data locally, but fetching ensures consistency
            fetchDrawings(currentProjectID);
        };

        eventSource.onerror = (err) => {
            // console.error("SSE Error:", err); 
            // Silent retry logic is handled by browser, but we close on fatal
            // eventSource.close();
        };

        return () => {
            if (eventSourceRef.current) {
                eventSourceRef.current.close();
            }
        };
    }, [currentProjectID]); // Re-run when project changes

    const handleAction = async (action, id) => {
        try {
            switch (action) {
                case 'claim': await drawingService.claim(id); break;
                case 'submit': await drawingService.submit(id); break;
                case 'release': await drawingService.release(id); break;
                case 'reject': await drawingService.reject(id); break;
            }
            // Optimistic update or wait for SSE
            fetchDrawings(currentProjectID);
        } catch (err) {
            alert(err.response?.data?.error || `${action} failed`);
        }
    };

    const myTasks = drawings.filter(d => d.assignee_id === parseInt(user.userID || user.user_id));
    const availableTasks = drawings.filter(d => d.assignee_id === null && d.current_stage !== 'approved');

    return (
        <div className="space-y-8 animate-in fade-in duration-500">
            <header className="flex justify-between items-end">
                <div>
                    <h1 className="text-4xl font-extrabold text-white">Dashboard</h1>
                    <p className="text-gray-400 mt-2">Manage and track quality control workflow.</p>
                </div>
                <ProjectSelector currentProject={currentProjectID} onSelect={setCurrentProjectID} />
            </header>

            {error && <div className="bg-red-900/20 border border-red-900/50 text-red-400 p-4 rounded-lg">{error}</div>}

            {loading && <div className="text-center py-10 text-blue-400">Loading workspace...</div>}

            {!loading && (
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                    {/* My Active Tasks */}
                    <section className="space-y-4">
                        <div className="flex items-center space-x-2">
                            <Clock className="text-blue-500" />
                            <h2 className="text-2xl font-bold">My Active Tasks</h2>
                        </div>
                        <div className="grid gap-4">
                            {myTasks.length === 0 ? (
                                <p className="text-gray-500 italic p-6 bg-gray-900/50 rounded-xl border border-gray-800">No active tasks assigned to you.</p>
                            ) : (
                                myTasks.map(drawing => (
                                    <TaskCard key={drawing.id} drawing={drawing} user={user} onAction={handleAction} />
                                ))
                            )}
                        </div>
                    </section>

                    {/* Available to Claim */}
                    <section className="space-y-4">
                        <div className="flex items-center space-x-2">
                            <UserPlus className="text-yellow-500" />
                            <h2 className="text-2xl font-bold">Available in Pool</h2>
                        </div>
                        <div className="grid gap-4">
                            {availableTasks.length === 0 ? (
                                <p className="text-gray-500 italic p-6 bg-gray-900/50 rounded-xl border border-gray-800">No drawings currently available for claim.</p>
                            ) : (
                                availableTasks.map(drawing => (
                                    <TaskCard key={drawing.id} drawing={drawing} user={user} onAction={handleAction} />
                                ))
                            )}
                        </div>
                    </section>
                </div>
            )}

            <OverviewTable drawings={drawings} />
        </div>
    );
};

export default Dashboard;
// End of Dashboard component
