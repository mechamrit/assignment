import React from 'react';

const OverviewTable = ({ drawings }) => {
    return (
        <section className="space-y-4 pt-4 border-t border-gray-800">
            <h2 className="text-2xl font-bold">All Drawings Overview</h2>
            <div className="overflow-x-auto rounded-xl border border-gray-800 shadow-2xl">
                <table className="w-full text-left bg-gray-900">
                    <thead className="bg-gray-800 text-gray-300 text-sm uppercase tracking-wider">
                        <tr>
                            <th className="px-6 py-4 font-semibold">Drawing Name</th>
                            <th className="px-6 py-4 font-semibold">Rev</th>
                            <th className="px-6 py-4 font-semibold">Stage</th>
                            <th className="px-6 py-4 font-semibold">Assigned To</th>
                            <th className="px-6 py-4 font-semibold">Last Updated</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-800">
                        {drawings.map(d => (
                            <tr key={d.id} className="hover:bg-gray-800/50 transition-colors">
                                <td className="px-6 py-4 font-medium">
                                    {d.title}
                                    <div className="text-xs text-gray-500">{d.description}</div>
                                </td>
                                <td className="px-6 py-4 font-mono text-sm text-purple-400">R{d.revision}</td>
                                <td className="px-6 py-4">
                                    <span className="px-3 py-1 bg-gray-800 rounded-full text-xs font-semibold capitalize border border-gray-700">
                                        {d.current_stage.replace('_', ' ')}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-400">
                                    {d.assignee ? d.assignee.username : '—'}
                                </td>
                                <td className="px-6 py-4 text-sm text-gray-500">
                                    {new Date(d.updated_at).toLocaleString()}
                                </td>
                            </tr>
                        ))}
                        {drawings.length === 0 && (
                            <tr>
                                <td colSpan="5" className="px-6 py-10 text-center text-gray-500 italic">No drawings found in system.</td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </section>
    );
};

export default OverviewTable;
