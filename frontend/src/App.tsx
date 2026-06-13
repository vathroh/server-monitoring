import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Dashboard from './pages/Dashboard';
import Login from './pages/Login';
import Profile from './pages/Profile';
import ServerList from './pages/ServerList';
import ServerForm from './pages/ServerForm';
import ServerDetail from './pages/ServerDetail';
import AlertList from './pages/AlertList';
import Settings from './pages/Settings';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/login" element={<Login />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/servers" element={<ServerList />} />
          <Route path="/servers/new" element={<ServerForm />} />
          <Route path="/servers/:id" element={<ServerDetail />} />
          <Route path="/servers/:id/edit" element={<ServerForm />} />
          <Route path="/alerts" element={<AlertList />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
