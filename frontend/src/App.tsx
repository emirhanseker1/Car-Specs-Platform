import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppShell from './components/AppShell';
import ScrollToTop from './components/ScrollToTop';
import Home from './pages/Home.tsx';
import ModelList from './pages/ModelList.tsx';
import ModelGenerations from './pages/ModelGenerations.tsx';
import PowertrainSelect from './pages/PowertrainSelect.tsx';
import VehicleDetails from './pages/VehicleDetails.tsx';
import Search from './pages/Search.tsx';
import GolfMkGroup from './pages/GolfMkGroup.tsx';
import Compare from './pages/Compare.tsx';
import Guides from './pages/Guides.tsx';
import TransmissionGuide from './pages/TransmissionGuide.tsx';
import EngineTermsGuide from './pages/EngineTermsGuide.tsx';
import About from './pages/About.tsx';

function App() {
  return (
    <Router>
      <ScrollToTop />
      <Routes>
        <Route element={<AppShell />}>
          <Route path="/" element={<Home />} />
          <Route path="/search" element={<Search />} />
          <Route path="/brand/:brandName" element={<ModelList />} />
          <Route path="/brand/:brandName/model/:modelName" element={<ModelGenerations />} />
          <Route path="/brand/:brandName/model/:modelName/mk/:mk" element={<GolfMkGroup />} />
          <Route path="/vehicle/:id/powertrain" element={<PowertrainSelect />} />
          <Route path="/vehicle/:id" element={<VehicleDetails />} />
          <Route path="/compare" element={<Compare />} />
          <Route path="/guides" element={<Guides />} />
          <Route path="/guides/transmission" element={<TransmissionGuide />} />
          <Route path="/guides/engine" element={<EngineTermsGuide />} />
          <Route path="/about" element={<About />} />

          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
