import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Box, Typography, Paper, Button, AppBar, Toolbar, CircularProgress } from '@mui/material';
import { serverService } from '../services/api';

export default function ServerDetail() {
  const { id } = useParams();
  const navigate = useNavigate();

  const { data, isLoading, error } = useQuery({
    queryKey: ['server', id],
    queryFn: () => serverService.get(Number(id)),
  });

  if (isLoading) return <Box p={4} display="flex" justifyContent="center"><CircularProgress /></Box>;
  if (error || !data?.success) return <Box p={4}>Error loading server details.</Box>;

  const server = data.data;

  return (
    <Box sx={{ flexGrow: 1, height: '100vh', backgroundColor: '#f0f2f5' }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Server Details
          </Typography>
          <Button color="inherit" onClick={() => navigate('/servers')}>
            Back to List
          </Button>
        </Toolbar>
      </AppBar>
      <Box p={4} display="flex" justifyContent="center">
        <Paper elevation={3} sx={{ p: 4, width: '100%', maxWidth: 600 }}>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
            <Typography variant="h5">{server.name}</Typography>
            <Button variant="outlined" onClick={() => navigate(`/servers/${server.id}/edit`)}>Edit</Button>
          </Box>
          <Typography variant="body1" gutterBottom><strong>Hostname:</strong> {server.hostname}</Typography>
          <Typography variant="body1" gutterBottom><strong>IP Address:</strong> {server.ip_address}</Typography>
          <Typography variant="body1" gutterBottom><strong>Environment:</strong> {server.environment || 'N/A'}</Typography>
          <Typography variant="body1" gutterBottom><strong>API Key:</strong> <span style={{ fontFamily: 'monospace', background: '#e0e0e0', padding: '2px 4px', borderRadius: '4px' }}>{server.api_key}</span></Typography>
          <Typography variant="body1" gutterBottom><strong>Description:</strong> {server.description || 'N/A'}</Typography>
          <Typography variant="body2" color="textSecondary" sx={{ mt: 3 }}>
            Added on {new Date(server.created_at).toLocaleString()}
          </Typography>
        </Paper>
      </Box>
    </Box>
  );
}
