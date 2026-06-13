import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, TablePagination, AppBar, Toolbar, Chip } from '@mui/material';
import { serverService } from '../services/api';

export default function ServerList() {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['servers', page + 1, rowsPerPage],
    queryFn: () => serverService.list(page + 1, rowsPerPage),
  });

  const deleteMutation = useMutation({
    mutationFn: serverService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
    },
  });

  const handleDelete = (id: number) => {
    if (confirm('Are you sure you want to delete this server?')) {
      deleteMutation.mutate(id);
    }
  };

  const servers = data?.data?.data || [];
  const totalCount = data?.data?.total || 0;

  return (
    <Box sx={{ flexGrow: 1, height: '100vh', backgroundColor: '#f0f2f5' }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Server Management
          </Typography>
          <Button color="inherit" onClick={() => navigate('/')}>
            Dashboard
          </Button>
        </Toolbar>
      </AppBar>
      <Box p={4}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h5">Servers</Typography>
          <Button variant="contained" color="primary" onClick={() => navigate('/servers/new')}>
            Add Server
          </Button>
        </Box>
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Hostname</TableCell>
                <TableCell>IP Address</TableCell>
                <TableCell>Environment</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">Loading...</TableCell>
                </TableRow>
              ) : servers.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} align="center">No servers found</TableCell>
                </TableRow>
              ) : (
                servers.map((server: any) => (
                  <TableRow key={server.id}>
                    <TableCell>{server.name}</TableCell>
                    <TableCell>{server.hostname}</TableCell>
                    <TableCell>{server.ip_address}</TableCell>
                    <TableCell>{server.environment}</TableCell>
                    <TableCell>
                      <Chip 
                        label={server.status || 'UNKNOWN'} 
                        size="small" 
                        color={
                          server.status === 'ONLINE' ? 'success' : 
                          server.status === 'WARNING' ? 'warning' : 'error'
                        } 
                      />
                    </TableCell>
                    <TableCell>
                      <Button size="small" onClick={() => navigate(`/servers/${server.id}`)}>View</Button>
                      <Button size="small" onClick={() => navigate(`/servers/${server.id}/edit`)}>Edit</Button>
                      <Button size="small" color="error" onClick={() => handleDelete(server.id)}>Delete</Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
          <TablePagination
            component="div"
            count={totalCount}
            page={page}
            onPageChange={(_, newPage) => setPage(newPage)}
            rowsPerPage={rowsPerPage}
            onRowsPerPageChange={(e) => {
              setRowsPerPage(parseInt(e.target.value, 10));
              setPage(0);
            }}
          />
        </TableContainer>
      </Box>
    </Box>
  );
}
