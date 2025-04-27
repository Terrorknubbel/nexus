import React, { useState } from "react";
import "./ProcessTable.scss";
import "bootstrap-icons/font/bootstrap-icons.css";
import { models } from "../../wailsjs/go/models";

type Props = {
  processes: models.Process[];
};

const stateIcons: Record<models.Process["state"], string> = {
  running: "🟢",
  sleeping: "🟡",
  stopped: "🔴",
};

type SortKey = keyof models.Process;
type SortDirection = "asc" | "desc";

const ProcessTable: React.FC<Props> = ({ processes }) => {
  const [sortKey, setSortKey] = useState<SortKey>("pid");
  const [sortDirection, setSortDirection] = useState<SortDirection>("asc");
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [isPaused, setIsPaused] = useState<boolean>(false);
  const [pinnedNames, setPinnedNames] = useState<Set<string>>(new Set());
  const [selectedPid, setSelectedPid] = useState<number | null>(null);
  const [expandedParents, setExpandedParents] = useState<Set<number>>(new Set());

  const toggleExpand = (pid: number) => {
    setExpandedParents((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(pid)) {
        newSet.delete(pid);
      } else {
        newSet.add(pid);
      }
      return newSet;
    });
  };


  const handleRowClick = (pid: number) => {
    setSelectedPid(pid === selectedPid ? null : pid);
  };

  const togglePin = (name: string) => {
    setPinnedNames((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(name)) {
        newSet.delete(name);
      } else {
        newSet.add(name);
      }
      return newSet;
    });
  };

  const isPinned = (name: string) => pinnedNames.has(name);

  const togglePause = () => {
    setIsPaused((prev) => !prev);
  };

  const handleSort = (key: SortKey) => {
    if (key === sortKey) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc");
    } else {
      setSortKey(key);
      setSortDirection("asc");
    }
  };

  const isSorted = (key: SortKey) => key === sortKey;

  const filteredProcesses = processes.filter((proc) => {
    const search = searchTerm.toLowerCase();
    return (
      proc.name.toLowerCase().includes(search) ||
      proc.state.toLowerCase().includes(search) ||
      proc.pid.toString().includes(search)
    );
  });

  const sortedProcesses = [...filteredProcesses].sort((a, b) => {
    const valA = a[sortKey];
    const valB = b[sortKey];

    if (valA == null && valB == null) return 0;
    if (valA == null) return 1;
    if (valB == null) return -1;

    if (valA < valB) return sortDirection === "asc" ? -1 : 1;
    if (valA > valB) return sortDirection === "asc" ? 1 : -1;
    return 0;
  });

  const getSortIcon = (key: SortKey) => {
    if (key !== sortKey) return <i className="bi bi-arrow-down-up"></i>;
    return sortDirection === "asc" ? (
      <i className="bi bi-arrow-up"></i>
    ) : (
      <i className="bi bi-arrow-down"></i>
    );
  };

  return (
    <div className="process-container">
        <div className="search-bar">
          <button
            className={`pause-button ${isPaused ? "active" : ""}`}
            onClick={togglePause}
            title={isPaused ? "Fortsetzen" : "Pausieren"}
          >
            <i className={`bi ${isPaused ? "bi-play-fill" : "bi-pause-fill"}`}></i>
          </button>
          <input
            type="text"
            placeholder="Suche nach Name, state oder PID..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <table className="process-table">
          <thead>
            <tr>
              <th onClick={() => handleSort("pid")} className={isSorted("pid") ? "sorted" : ""}>
                PID {getSortIcon("pid")}
              </th>
              <th onClick={() => handleSort("name")} className={isSorted("name") ? "sorted" : ""}>
                Name {getSortIcon("name")}
              </th>
              <th onClick={() => handleSort("state")} className={isSorted("state") ? "sorted" : ""}>
                state {getSortIcon("state")}
              </th>
            </tr>
          </thead>
            <tbody>
              {sortedProcesses.map((proc) => {
                const isChild    = proc.parent_pid > 0;
                const isVisible  = !isChild || expandedParents.has(proc.parent_pid);


                if (!isVisible) return null;

                const isParent = processes.some(p => p.parent_pid === proc.pid);
                return (
                  <tr
                    key={proc.pid}
                    className={`${proc.pid === selectedPid ? "selected" : ""} ${isChild ? "child" : ""}`}
                    onClick={() => handleRowClick(proc.pid)}
                  >
                    <td>
                      {isParent && (
                        <button
                          className="expand-button"
                          onClick={(e) => {
                            e.stopPropagation();
                            toggleExpand(proc.pid);
                          }}
                        >
                          {expandedParents.has(proc.pid) ? "🔽" : "▶️"}
                        </button>
                      )}
                      {!isParent && isChild && <span className="child-indent">↳</span>}
                      {proc.pid}
                    </td>
                    <td>{proc.name}</td>
                    <td className={`state ${proc.state}`}>
                      <span className="state-icon">{stateIcons[proc.state]}</span> {proc.state}
                    </td>
                  </tr>
                );
              })}
            </tbody>
        </table>
        <div className="action-bar">
          <button
            className="terminate-button"
            disabled={selectedPid === null}
            onClick={() => {
              if (selectedPid !== null) {
                console.log("Beende Prozess mit PID:", selectedPid);
                // Hier könnte dein Kill-Call kommen
              }
            }}
          >
            Prozess beenden
          </button>
        </div>
    </div>
  );
};

export default ProcessTable;
