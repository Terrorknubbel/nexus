import { ThemeProvider } from "styled-components";
import { GetAllProcesses } from "../wailsjs/go/main/App";
import { models } from "../wailsjs/go/models";
import ProcessTable from "./components/ProcessTable";
import Searchbar from "./components/Searchbar";
import { useEffect, useState } from 'react';
import theme from "./theme/theme";

function App() {
  const [processes, setProcesses] = useState<models.Process[]>([]);
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Funktion zum Abrufen der Prozesse
    const fetchProcesses = () => {
      GetAllProcesses()
        .then((data) => {
          setProcesses(data);
          setLoading(false);
        })
        .catch((err) => {
          console.error("Fehler beim Laden der Prozesse:", err);
        });
    };

    // Sofort einmal abrufen
    fetchProcesses();
    // Intervall setzen (5000ms = 5s)
    const intervalId = setInterval(fetchProcesses, 5000);

    // Cleanup: Intervall stoppen, wenn Komponente unmountet
    return () => clearInterval(intervalId);
  }, []);

  if (loading) {
    return <>Loading processes…</>;
  }

  return (
    <ThemeProvider theme={theme}>
      <Searchbar searchTerm={searchTerm} onSearch={setSearchTerm} />
      <ProcessTable searchTerm={searchTerm} processes={processes} />
    </ThemeProvider>
  );
}

export default App;

