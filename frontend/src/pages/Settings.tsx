import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Typography, Paper, AppBar, Toolbar, TextField, Switch, FormControlLabel, Select, MenuItem, InputLabel, FormControl, Divider } from '@mui/material';
import { settingService } from '../services/api';

export default function Settings() {
  const navigate = useNavigate();
  const [telegramEnabled, setTelegramEnabled] = useState(false);
  const [telegramBotToken, setTelegramBotToken] = useState('');
  const [telegramChatId, setTelegramChatId] = useState('');
  
  const [whatsappEnabled, setWhatsappEnabled] = useState(false);
  const [whatsappProvider, setWhatsappProvider] = useState('waha');
  const [whatsappEndpoint, setWhatsappEndpoint] = useState('');
  const [whatsappChatId, setWhatsappChatId] = useState('');

  const [saving, setSaving] = useState(false);

  useEffect(() => {
    settingService.getAll().then(res => {
      const data = res.data || {};
      setTelegramEnabled(data['telegram_enabled'] === 'true');
      setTelegramBotToken(data['telegram_bot_token'] || '');
      setTelegramChatId(data['telegram_chat_id'] || '');

      setWhatsappEnabled(data['whatsapp_enabled'] === 'true');
      setWhatsappProvider(data['whatsapp_provider'] || 'waha');
      setWhatsappEndpoint(data['whatsapp_endpoint'] || '');
      setWhatsappChatId(data['whatsapp_chat_id'] || '');
    });
  }, []);

  const handleSave = async () => {
    setSaving(true);
    await settingService.save({
      telegram_enabled: telegramEnabled ? 'true' : 'false',
      telegram_bot_token: telegramBotToken,
      telegram_chat_id: telegramChatId,
      whatsapp_enabled: whatsappEnabled ? 'true' : 'false',
      whatsapp_provider: whatsappProvider,
      whatsapp_endpoint: whatsappEndpoint,
      whatsapp_chat_id: whatsappChatId,
    });
    setSaving(false);
    alert('Settings saved successfully!');
  };

  return (
    <Box sx={{ flexGrow: 1, minHeight: '100vh', backgroundColor: '#f0f2f5' }}>
      <AppBar position="static" elevation={0} sx={{ backgroundColor: '#1976d2' }}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
            Velocity Monitoring
          </Typography>
          <Button color="inherit" onClick={() => navigate('/')}>Dashboard</Button>
          <Button color="inherit" onClick={() => navigate('/servers')}>Servers</Button>
          <Button color="inherit" onClick={() => navigate('/alerts')}>Alerts</Button>
        </Toolbar>
      </AppBar>
      <Box p={4} maxWidth="800px" margin="0 auto">
        <Typography variant="h5" fontWeight="bold" mb={3}>Settings</Typography>
        <Paper sx={{ p: 4 }}>
          <Typography variant="h6" mb={2}>Telegram Notifications</Typography>
          <FormControlLabel
            control={<Switch checked={telegramEnabled} onChange={(e) => setTelegramEnabled(e.target.checked)} />}
            label="Enable Telegram Alerts"
          />
          <Box mt={2}>
            <TextField
              fullWidth
              label="Bot Token"
              variant="outlined"
              margin="normal"
              value={telegramBotToken}
              onChange={(e) => setTelegramBotToken(e.target.value)}
              disabled={!telegramEnabled}
            />
            <TextField
              fullWidth
              label="Chat ID"
              variant="outlined"
              margin="normal"
              value={telegramChatId}
              onChange={(e) => setTelegramChatId(e.target.value)}
              disabled={!telegramEnabled}
            />
          </Box>

          <Divider sx={{ my: 4 }} />

          <Typography variant="h6" mb={2}>WhatsApp Notifications</Typography>
          <FormControlLabel
            control={<Switch checked={whatsappEnabled} onChange={(e) => setWhatsappEnabled(e.target.checked)} />}
            label="Enable WhatsApp Alerts"
          />
          <Box mt={2}>
            <FormControl fullWidth margin="normal" disabled={!whatsappEnabled}>
              <InputLabel>Provider</InputLabel>
              <Select
                value={whatsappProvider}
                label="Provider"
                onChange={(e) => setWhatsappProvider(e.target.value)}
              >
                <MenuItem value="waha">WAHA (WhatsApp HTTP API)</MenuItem>
                <MenuItem value="evolution">Evolution API</MenuItem>
                <MenuItem value="meta">Meta Cloud API</MenuItem>
              </Select>
            </FormControl>
            <TextField
              fullWidth
              label="Endpoint URL (e.g. http://localhost:3000/api/sendText)"
              variant="outlined"
              margin="normal"
              value={whatsappEndpoint}
              onChange={(e) => setWhatsappEndpoint(e.target.value)}
              disabled={!whatsappEnabled}
            />
            <TextField
              fullWidth
              label="Chat ID (e.g. 123456789@c.us)"
              variant="outlined"
              margin="normal"
              value={whatsappChatId}
              onChange={(e) => setWhatsappChatId(e.target.value)}
              disabled={!whatsappEnabled}
            />
          </Box>

          <Box mt={3} display="flex" justifyContent="flex-end">
            <Button variant="contained" onClick={handleSave} disabled={saving}>
              {saving ? 'Saving...' : 'Save Settings'}
            </Button>
          </Box>
        </Paper>
      </Box>
    </Box>
  );
}
