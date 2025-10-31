import React, { useState } from "react";
import {
  Box,
  Container,
  Typography,
  Paper,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
} from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";

interface User {
  id: number;
  name: string;
  role: "admin" | "user";
}

const AdminDashboard: React.FC = () => {
  const [users] = useState<User[]>([
    { id: 1, name: "Jane Doe", role: "admin" },
    { id: 2, name: "John Smith", role: "user" },
    { id: 3, name: "Emily Johnson", role: "user" },
  ]);

  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.clear();
    navigate("/login");
  };

  return (
    <Box
      sx={{
        minHeight: "100vh",
        background: "linear-gradient(to right, #dbeafe, #e0f2fe, #f0fdfa)",
        py: 6,
      }}
    >
      <Container maxWidth="md">
        <Paper elevation={8} sx={{ p: 4, borderRadius: 3 }}>
          {/* Header */}
          <Box
            display="flex"
            justifyContent="space-between"
            alignItems="center"
            mb={4}
          >
            <Typography variant="h4" fontWeight="bold" color="primary">
              Admin Dashboard
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
              onClick={handleLogout}
            >
              Logout
            </Button>
          </Box>

          {/* Users Table */}
          <Typography variant="h6" fontWeight="medium" mb={2}>
            All Users
          </Typography>

          <TableContainer component={Paper} elevation={3}>
            <Table>
              <TableHead>
                <TableRow sx={{ backgroundColor: "#f3f4f6" }}>
                  <TableCell>ID</TableCell>
                  <TableCell>Name</TableCell>
                  <TableCell>Role</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {users.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell>{user.id}</TableCell>
                    <TableCell>{user.name}</TableCell>
                    <TableCell>
                      <Chip
                        label={user.role.toUpperCase()}
                        color={user.role === "admin" ? "primary" : "success"}
                        size="small"
                      />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      </Container>
    </Box>
  );
};

export default AdminDashboard;
