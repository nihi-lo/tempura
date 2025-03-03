"use client";

import React from "react";

interface HelloButtonProps {
  message: string;
}

export const HelloButton: React.FC<HelloButtonProps> = ({ message }) => {
  const handleButtonClick = () => {
    console.log(message);
  };

  return (
    <button
      className="h-10 min-w-20 rounded bg-blue-500 px-4 text-sm text-white hover:bg-blue-700"
      onClick={handleButtonClick}
    >
      Click!
    </button>
  );
};
