import { Box, Typography, Paper, Button, AppBar, Toolbar } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';

export default function Profile() {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);

  return (
    <Box sx={{ flexGrow: 1, height: '100vh', backgroundColor: '#f0f2f5' }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            User Profile
          </Typography>
          <Button color="inherit" onClick={() => navigate('/')}>
            Back to Dashboard
          </Button>
        </Toolbar>
      </AppBar>
      <Box p={4} display="flex" justifyContent="center">
        <Paper elevation={3} sx={{ p: 4, width: '100%', maxWidth: 600 }}>
          <Typography variant="h5" gutterBottom>Profile Details</Typography>
          <Typography variant="body1" sx={{ mt: 2 }}>
            <strong>Email:</strong> {user?.email || 'Unknown'}
          </Typography>
          <Typography variant="body1" sx={{ mt: 1 }}>
            <strong>Joined:</strong> {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'Unknown'}
          </Typography>
        </Paper>
      </Box>
    </Box>
  );
}
