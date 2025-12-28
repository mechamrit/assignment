import React from 'react';
import { ArrowRight, RotateCcw, XCircle, CheckCircle, AlertCircle, Clock } from 'lucide-react';

const TaskCard = ({ drawing, user, onAction }) => {
    const getStageIcon = (stage) => {
        switch (stage) {
            case 'approved': return <CheckCircle className="w-4 h-4 text-green-500" />;
            case 'unassigned': return <AlertCircle className="w-4 h-4 text-yellow-500" />;
            default: return <Clock className="w-4 h-4 text-blue-500" />;
        }
    };

    const isAssignedToMe = drawing.assignee_id === parseInt(user.userID || user.user_id);
    const inQC = drawing.current_stage === 'first_qc' || drawing.current_stage === 'final_qc';

    return (
        <div className="bg-gray-900 border border-gray-800 p-5 rounded-xl shadow-lg hover:border-blue-500/50 transition-all flex justify-between items-center group">
            <div>
                <div className="flex items-center space-x-2">
                    {getStageIcon(drawing.current_stage)}
                    <h3 className="font-bold text-lg">{drawing.title}</h3>
                    <span className="text-[10px] bg-purple-900/30 text-purple-400 px-1.5 py-0.5 rounded border border-purple-800/50">REV {drawing.revision}</span>
                </div>
                <div className="flex items-center space-x-3 mt-1 text-sm">
                    <span className="px-2 py-0.5 bg-blue-900/30 text-blue-400 rounded-md border border-blue-800/50 capitalize font-mono text-xs">
                        {drawing.current_stage.replace('_', ' ')}
                    </span>
                    {drawing.description && (
                        <span className="text-gray-500 text-xs truncate max-w-[200px]">{drawing.description}</span>
                    )}
                </div>
            </div>

            <div className="flex items-center space-x-3">
                {isAssignedToMe ? (
                    <>
                        {inQC && (
                            <button
                                onClick={() => onAction('reject', drawing.id)}
                                className="flex items-center space-x-1.5 bg-red-900/40 hover:bg-red-800 text-red-400 px-3 py-2 rounded-lg text-sm font-bold border border-red-800/50 transition-all"
                                title="Send back for rework"
                            >
                                <RotateCcw className="w-4 h-4" />
                                <span>Reject</span>
                            </button>
                        )}
                        <button
                            onClick={() => onAction('release', drawing.id)}
                            className="flex items-center space-x-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 px-3 py-2 rounded-lg text-sm font-bold border border-gray-700 transition-all"
                            title="Release back to pool"
                        >
                            <XCircle className="w-4 h-4" />
                            <span>Release</span>
                        </button>
                        <button
                            onClick={() => onAction('submit', drawing.id)}
                            className="flex items-center space-x-2 bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg font-bold transition-transform group-hover:scale-105 shadow-lg"
                        >
                            <span>Submit</span>
                            <ArrowRight className="w-4 h-4" />
                        </button>
                    </>
                ) : (
                    <button
                        onClick={() => onAction('claim', drawing.id)}
                        className="bg-yellow-600 hover:bg-yellow-700 text-white px-4 py-2 rounded-lg font-bold transition-all"
                    >
                        Claim
                    </button>
                )}
            </div>
        </div>
    );
};

export default TaskCard;
