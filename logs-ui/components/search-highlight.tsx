import React from 'react';

interface SearchHighlightProps {
  text: string;
  searchQuery: string;
  className?: string;
}

/**
 * Component to highlight search matches in text
 */
export function SearchHighlight({ text, searchQuery, className = '' }: SearchHighlightProps) {
  if (!searchQuery || !text) {
    return <span className={className}>{text}</span>;
  }

  const parts = text.split(new RegExp(`(${searchQuery})`, 'gi'));

  return (
    <span className={className}>
      {parts.map((part, index) => 
        part.toLowerCase() === searchQuery.toLowerCase() ? (
          <mark key={index} className="bg-yellow-500/30 text-yellow-200 px-0.5 rounded">
            {part}
          </mark>
        ) : (
          <span key={index}>{part}</span>
        )
      )}
    </span>
  );
}


