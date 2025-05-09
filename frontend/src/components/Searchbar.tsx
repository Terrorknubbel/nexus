import React, { useState } from "react";
import styled from "styled-components";

type Props = {
  searchTerm: string;
  onSearch: (searchTerm: string) => void;
};

const PauseButton = styled.button<{ isPaused: boolean }>`
  width: 2.5rem;
  height: 2.5rem;
  border: 1px solid ${(props) => props.theme.primary};
  background-color: ${(props) =>
    props.isPaused ? props.theme.primary : props.theme.surface
  };
  border-radius: 0.5rem;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.3s ease;
  i {
    font-size: 1.2rem;
  }
  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
`;

const Searchbar: React.FC<Props> = ({ searchTerm, onSearch }) => {
  const [isPaused, setIsPaused] = useState<boolean>(false);

  const togglePause = () => {
    setIsPaused((prev) => !prev);
  };

  return (
    <div className="search-bar">
      <PauseButton
        isPaused={isPaused}
        className={`pause-button}`}
        onClick={togglePause}
        title={isPaused ? "Fortsetzen" : "Pausieren"}
      >
        <i className={`bi ${isPaused ? "bi-play-fill" : "bi-pause-fill"}`}></i>
      </PauseButton>
      <input
        type="text"
        placeholder="Suche nach Name, state oder PID..."
        value={searchTerm}
        onChange={(e) => onSearch(e.target.value)}
      />
    </div>
  );
}

export default Searchbar;
