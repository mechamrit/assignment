import React, { useEffect, useState } from 'react';
import { Layers } from 'lucide-react';

const ProjectSelector = ({ currentProject, onSelect }) => {
    // In a real app, we'd fetch these from /api/v1/projects
    // For now, we'll hardcode the ones we seeded.
    const projects = [
        { id: 1, name: "Aerodynamics Suite" },
        { id: 2, name: "Power Unit" }
    ];

    return (
        <div className="flex items-center space-x-4 bg-gray-800 p-2 rounded-lg border border-gray-700">
            <Layers className="text-gray-400" />
            <select
                className="bg-transparent text-white font-bold outline-none cursor-pointer"
                value={currentProject}
                onChange={(e) => onSelect(parseInt(e.target.value))}
            >
                {projects.map(p => (
                    <option key={p.id} value={p.id} className="bg-gray-900 text-white">
                        {p.name}
                    </option>
                ))}
            </select>
        </div>
    );
};

export default ProjectSelector;
