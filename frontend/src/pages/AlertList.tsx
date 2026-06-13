import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, AppBar, Toolbar, Chip, ToggleButtonGroup, ToggleButton } from '@mui/material';
import { alertService } from '../services/api';

export default function AlertList() {
  const navigate = useNavigate();
  const [stateFilter, setStateFilter] = useState<string>('OPEN');

  const { data, isLoading } = useQuery({
    queryKey: ['alerts', stateFilter],
    queryFn: () => alertService.list(stateFilter === 'ALL' ? undefined : stateFilter),
    refetchInterval: 10000,
  });

  const alerts = data?.data || [];

  return (
    <Box sx={{ flexGrow: 1, minHeight: '100vh', backgroundColor: '#f0f2f5' }}>
      <AppBar position="static" elevation={0} sx={{ backgroundColor: '#1976d2' }}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
            Velocity Monitoring
          </Typography>
          <Button color="inherit" onClick={() => navigate('/')}>Dashboard</Button>
          <Button color="inherit" onClick={() => navigate('/servers')}>Servers</Button>
        </Toolbar>
      </AppBar>
      <Box p={4} maxWidth="1200px" margin="0 auto">
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h5" fontWeight="bold">Alerts</Typography>
          <ToggleButtonGroup
            size="small"
            value={stateFilter}
            exclusive
            onChange={(_, val) => val && setStateFilter(val)}
          >
            <ToggleButton value="ALL">All</ToggleButton>
            <ToggleButton value="OPEN">Open</ToggleButton>
            <ToggleButton value="RESOLVED">Resolved</ToggleButton>
          </ToggleButtonGroup>
        </Box>
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Server</TableCell>
                <TableCell>Rule</TableCell>
                <TableCell>Severity</TableCell>
                <TableCell>State</TableCell>
                <TableCell>Message</TableCell>
                <TableCell>Created At</TableCell>
                <TableCell>Resolved At</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {isLoading ? (
                <TableRow><TableCell colSpan={7} align="center">Loading...</TableCell></TableRow>
              ) : alerts.length === 0 ? (
                <TableRow><TableCell colSpan={7} align="center">No alerts found</TableCell></TableRow>
              ) : (
                alerts.map((alert: any) => (
                  <TableRow key={alert.id}>
                    <TableCell>{alert.server?.name || `ID ${alert.server_id}`}</TableCell>
                    <TableCell><Chip label={alert.rule} size="small" /></TableCell>
                    <TableCell>
                      <Chip 
                        label={alert.severity} 
                        size="small" 
                        color={alert.severity === 'CRITICAL' ? 'error' : 'warning'} 
                      />
                    </TableCell>
                    <TableCell>
                      <Chip 
                        label={alert.state} 
                        size="small" 
                        color={alert.state === 'OPEN' ? 'error' : 'success'} 
                        variant={alert.state === 'OPEN' ? 'filled' : 'outlined'}
                      />
                    </TableCell>
                    <TableCell>{alert.message}</TableCell>
                    <TableCell>{new Date(alert.created_at).toLocaleString()}</TableCell>
                    <TableCell>{alert.resolved_at ? new Date(alert.resolved_at).toLocaleString() : '-'}</TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>
    </Box>
  );
}
