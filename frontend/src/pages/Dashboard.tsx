import { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  Box,
  Typography,
  Grid,
  Paper,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  ToggleButtonGroup,
  ToggleButton,
  AppBar,
  Toolbar,
  Button,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { dashboardService, serverService } from "../services/api";
import { useAuthStore } from "../store/authStore";

export default function Dashboard() {
  const navigate = useNavigate();
  const logout = useAuthStore((state) => state.logout);

  const [selectedServerId, setSelectedServerId] = useState<number | "">("");
  const [timeRange, setTimeRange] = useState<string>("1h");

  // Fetch Dashboard Summary
  const { data: summaryData } = useQuery({
    queryKey: ["dashboard", "summary"],
    queryFn: dashboardService.getSummary,
    refetchInterval: 2000,
  });
  const summary = summaryData?.data || {
    total_servers: 0,
    online_servers: 0,
    warning_servers: 0,
    offline_servers: 0,
  };

  // Fetch all servers for the dropdown selector
  const { data: serversData } = useQuery({
    queryKey: ["servers", "all"],
    queryFn: () => serverService.list(1, 100),
  });
  const servers = serversData?.data?.data || [];

  // Default to first server if none selected
  useEffect(() => {
    if (selectedServerId === "" && servers.length > 0) {
      setSelectedServerId(servers[0].id);
    }
  }, [servers, selectedServerId]);

  // Fetch Metrics Trend for selected server
  const { data: trendData } = useQuery({
    queryKey: ["dashboard", "trend", selectedServerId, timeRange],
    queryFn: () =>
      dashboardService.getTrend(Number(selectedServerId), timeRange),
    enabled: selectedServerId !== "",
    refetchInterval: 2000,
  });
  const chartData = trendData?.data || [];

  // Format data for charts
  const formattedChartData = chartData.map((d: any) => ({
    time: new Date(d.created_at).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
    }),
    cpu: parseFloat(d.cpu_usage.toFixed(2)),
    memory: parseFloat(d.memory_usage.toFixed(2)),
    disk: parseFloat(d.disk_usage.toFixed(2)),
  }));

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <Box sx={{ flexGrow: 1, minHeight: "100vh", backgroundColor: "#f4f6f8" }}>
      <AppBar
        position="static"
        elevation={0}
        sx={{ backgroundColor: "#1976d2" }}
      >
        <Toolbar>
          <Typography
            variant="h6"
            component="div"
            sx={{ flexGrow: 1, fontWeight: "bold" }}
          >
            Velocity Monitoring
          </Typography>
          <Button color="inherit" onClick={() => navigate("/servers")}>
            Servers
          </Button>
          <Button color="inherit" onClick={() => navigate("/alerts")}>
            Alerts
          </Button>
          <Button color="inherit" onClick={() => navigate("/settings")}>
            Settings
          </Button>
          <Button color="inherit" onClick={() => navigate("/profile")}>
            Profile
          </Button>
          <Button color="inherit" onClick={handleLogout}>
            Logout
          </Button>
        </Toolbar>
      </AppBar>

      <Box p={4} maxWidth="1200px" margin="0 auto">
        <Typography variant="h5" mb={3} fontWeight="bold" color="text.primary">
          Dashboard Summary
        </Typography>

        {/* Widget Row */}
        <Grid container spacing={3} mb={4}>
          {[
            {
              title: "Total Servers",
              value: summary.total_servers,
              color: "#1976d2",
            },
            {
              title: "Online Servers",
              value: summary.online_servers,
              color: "#2e7d32",
            },
            {
              title: "Warning Servers",
              value: summary.warning_servers,
              color: "#ed6c02",
            },
            {
              title: "Active Alerts",
              value: summary.active_alerts,
              color: "#d32f2f",
            },
          ].map((widget, i) => (
            <Grid item xs={12} sm={6} md={3} key={i}>
              <Paper
                elevation={2}
                sx={{ p: 3, textAlign: "center", borderRadius: 2 }}
              >
                <Typography variant="subtitle2" color="textSecondary">
                  {widget.title}
                </Typography>
                <Typography
                  variant="h4"
                  sx={{ color: widget.color, fontWeight: "bold", mt: 1 }}
                >
                  {widget.value}
                </Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>

        {/* Chart Controls */}
        <Box
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          mb={2}
        >
          <Typography variant="h6" fontWeight="bold">
            Server Metrics
          </Typography>
          <Box display="flex" gap={2}>
            <FormControl size="small" sx={{ minWidth: 200 }}>
              <InputLabel>Server</InputLabel>
              <Select
                value={selectedServerId}
                label="Server"
                onChange={(e) => setSelectedServerId(e.target.value as number)}
              >
                {servers.map((s: any) => (
                  <MenuItem key={s.id} value={s.id}>
                    {s.name} ({s.hostname})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <ToggleButtonGroup
              size="small"
              value={timeRange}
              exclusive
              onChange={(_, val) => val && setTimeRange(val)}
            >
              <ToggleButton value="1h">1H</ToggleButton>
              <ToggleButton value="24h">24H</ToggleButton>
              <ToggleButton value="7d">7D</ToggleButton>
            </ToggleButtonGroup>
          </Box>
        </Box>

        {/* Charts */}
        <Grid container spacing={3}>
          {/* CPU Chart */}
          <Grid item xs={12}>
            <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
              <Typography variant="subtitle1" mb={2} fontWeight="bold">
                CPU Usage (%)
              </Typography>
              <Box height={250}>
                {formattedChartData.length > 0 ? (
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={formattedChartData}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} />
                      <XAxis dataKey="time" tick={{ fontSize: 12 }} />
                      <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
                      <Tooltip />
                      <Line
                        type="monotone"
                        dataKey="cpu"
                        stroke="#1976d2"
                        strokeWidth={2}
                        dot={false}
                        isAnimationActive={false}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                ) : (
                  <Box
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    height="100%"
                  >
                    <Typography color="textSecondary">
                      No data available.
                    </Typography>
                  </Box>
                )}
              </Box>
            </Paper>
          </Grid>

          {/* Memory Chart */}
          <Grid item xs={12} md={6}>
            <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
              <Typography variant="subtitle1" mb={2} fontWeight="bold">
                Memory Usage (%)
              </Typography>
              <Box height={250}>
                {formattedChartData.length > 0 ? (
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={formattedChartData}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} />
                      <XAxis dataKey="time" tick={{ fontSize: 12 }} />
                      <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
                      <Tooltip />
                      <Line
                        type="monotone"
                        dataKey="memory"
                        stroke="#9c27b0"
                        strokeWidth={2}
                        dot={false}
                        isAnimationActive={false}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                ) : (
                  <Box
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    height="100%"
                  >
                    <Typography color="textSecondary">
                      No data available.
                    </Typography>
                  </Box>
                )}
              </Box>
            </Paper>
          </Grid>

          {/* Disk Chart */}
          <Grid item xs={12} md={6}>
            <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
              <Typography variant="subtitle1" mb={2} fontWeight="bold">
                Disk Usage (%)
              </Typography>
              <Box height={250}>
                {formattedChartData.length > 0 ? (
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={formattedChartData}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} />
                      <XAxis dataKey="time" tick={{ fontSize: 12 }} />
                      <YAxis domain={[0, 100]} tick={{ fontSize: 12 }} />
                      <Tooltip />
                      <Line
                        type="monotone"
                        dataKey="disk"
                        stroke="#0288d1"
                        strokeWidth={2}
                        dot={false}
                        isAnimationActive={false}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                ) : (
                  <Box
                    display="flex"
                    alignItems="center"
                    justifyContent="center"
                    height="100%"
                  >
                    <Typography color="textSecondary">
                      No data available.
                    </Typography>
                  </Box>
                )}
              </Box>
            </Paper>
          </Grid>
        </Grid>
      </Box>
    </Box>
  );
}
