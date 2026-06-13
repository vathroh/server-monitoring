import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Box, Button, Typography, Paper, TextField, AppBar, Toolbar } from '@mui/material';
import { serverService } from '../services/api';

export default function ServerForm() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [formData, setFormData] = useState({
    name: '',
    hostname: '',
    ip_address: '',
    environment: '',
    description: '',
  });

  const { data, isLoading } = useQuery({
    queryKey: ['server', id],
    queryFn: () => serverService.get(Number(id)),
    enabled: isEdit,
  });

  useEffect(() => {
    if (isEdit && data?.success) {
      const s = data.data;
      setFormData({
        name: s.name,
        hostname: s.hostname,
        ip_address: s.ip_address,
        environment: s.environment,
        description: s.description,
      });
    }
  }, [isEdit, data]);

  const saveMutation = useMutation({
    mutationFn: (payload: any) => isEdit ? serverService.update(Number(id), payload) : serverService.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['servers'] });
      navigate('/servers');
    },
    onError: (error: any) => {
      alert(error.response?.data?.message || 'Failed to save server');
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    saveMutation.mutate(formData);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  if (isEdit && isLoading) return <Box p={4}>Loading...</Box>;

  return (
    <Box sx={{ flexGrow: 1, minHeight: '100vh', backgroundColor: '#f0f2f5' }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            {isEdit ? 'Edit Server' : 'Add New Server'}
          </Typography>
          <Button color="inherit" onClick={() => navigate('/servers')}>
            Cancel
          </Button>
        </Toolbar>
      </AppBar>
      <Box p={4} display="flex" justifyContent="center">
        <Paper elevation={3} sx={{ p: 4, width: '100%', maxWidth: 600 }}>
          <Typography variant="h5" gutterBottom>{isEdit ? 'Edit Server' : 'Add Server'}</Typography>
          <form onSubmit={handleSubmit}>
            <TextField
              fullWidth label="Name" name="name" margin="normal"
              value={formData.name} onChange={handleChange} required
            />
            <TextField
              fullWidth label="Hostname" name="hostname" margin="normal"
              value={formData.hostname} onChange={handleChange} required
            />
            <TextField
              fullWidth label="IP Address" name="ip_address" margin="normal"
              value={formData.ip_address} onChange={handleChange} required
            />
            <TextField
              fullWidth label="Environment" name="environment" margin="normal"
              value={formData.environment} onChange={handleChange}
              placeholder="e.g. Production, Staging"
            />
            <TextField
              fullWidth label="Description" name="description" margin="normal"
              value={formData.description} onChange={handleChange}
              multiline rows={4}
            />
            <Box mt={3} display="flex" justifyContent="flex-end" gap={2}>
              <Button variant="outlined" onClick={() => navigate('/servers')}>Cancel</Button>
              <Button type="submit" variant="contained" color="primary" disabled={saveMutation.isPending}>
                {saveMutation.isPending ? 'Saving...' : 'Save'}
              </Button>
            </Box>
          </form>
        </Paper>
      </Box>
    </Box>
  );
}
