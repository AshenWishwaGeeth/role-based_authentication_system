import React, { useState } from "react";
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  Divider,
  Stack,
} from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";

interface User {
  name: string;
  email: string;
}

const UserDashboard: React.FC = () => {
  const [user] = useState<User>({ name: "John Doe", email: "john@example.com" });

  return (
    <Box
      sx={{
        minHeight: "100vh",
        background: "linear-gradient(to right, #bbf7d0, #e0f2fe, #a5b4fc)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        py: 6,
      }}
    >
      <Container maxWidth="sm">
        <Paper
          elevation={8}
          sx={{
            p: 4,
            borderRadius: 3,
            backdropFilter: "blur(10px)",
          }}
        >
          {/* Header */}
          <Stack direction="row" justifyContent="space-between" alignItems="center" mb={3}>
            <Typography variant="h4" fontWeight="bold" color="success.main">
              User Dashboard
            </Typography>
            <Button
              variant="contained"
              color="error"
              startIcon={<LogoutIcon />}
              sx={{
                textTransform: "none",
                fontWeight: "bold",
                "&:hover": { transform: "scale(1.05)" },
              }}
            >
              Logout
            </Button>
          </Stack>

          <Typography variant="body1" mb={3}>
            Welcome back, <strong>{user.name}</strong>!
          </Typography>

          <Divider sx={{ mb: 3 }} />

          {/* Profile Card */}
          <Paper elevation={2} sx={{ p: 3, borderRadius: 2 }}>
            <Typography variant="h6" fontWeight="medium" mb={2}>
              Profile Information
            </Typography>
            <Stack spacing={1}>
              <Typography variant="body2" color="text.secondary">
                <strong>Name:</strong> {user.name}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                <strong>Email:</strong> {user.email}
              </Typography>
            </Stack>
          </Paper>
        </Paper>
      </Container>
    </Box>
  );
};

export default UserDashboard;
