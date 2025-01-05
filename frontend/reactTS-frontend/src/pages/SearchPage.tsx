import React, { useState } from 'react';
import { Thread } from '../models/models';

const SearchPage: React.FC = () => {
    const [searchTerm, setSearchTerm] = useState<string>('');
    const [results, setResults] = useState<Thread[]>([]);

    const handleSearch = async () => {
        // Replace with your actual search logic, e.g., API call
        const response = await fetch(`/api/threads?search=${searchTerm}`);
        const data: Thread[] = await response.json();
        setResults(data);
    };

    return (
        <div>
            <h1>Search Threads</h1>
            <input
                type="text"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                placeholder="Enter keywords"
            />
            <button onClick={handleSearch}>Search</button>
            <div>
                {results.length > 0 ? (
                    <ul>
                        {results.map((thread) => (
                            <li key={thread.id}>{thread.title}</li>
                        ))}
                    </ul>
                ) : (
                    <p>No results found</p>
                )}
            </div>
        </div>
    );
};

export default SearchPage