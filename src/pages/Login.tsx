import React, { useState } from "react";
import {
  Avatar,
  Box,
  Button,
  Container,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Typography,
  Paper,
} from "@mui/material";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";

const Login: React.FC = () => {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [role, setRole] = useState<string>("user");

  return (
    <Box
      sx={{
        minHeight: "100vh",
        background: "linear-gradient(to bottom right, #a5b4fc, #c7d2fe, #bbf7d0)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Container maxWidth="xs">
        <Paper
          elevation={6}
          sx={{
            p: 4,
            borderRadius: 4,
            textAlign: "center",
            backdropFilter: "blur(10px)",
          }}
        >
          {/* Icon and Title */}
          <Avatar sx={{ bgcolor: "primary.main", mx: "auto", mb: 2 }}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography
            variant="h4"
            fontWeight="bold"
            color="primary"
            sx={{ mb: 1 }}
          >
            Welcome Back
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
            Sign in to your account
          </Typography>

          {/* Form */}
          <Box
            component="form"
            sx={{ display: "flex", flexDirection: "column", gap: 3 }}
          >
            <TextField
              label="Email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              fullWidth
              required
              variant="outlined"
            />

            <TextField
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              fullWidth
              required
              variant="outlined"
            />

            <FormControl fullWidth>
              <InputLabel>Role</InputLabel>
              <Select
                value={role}
                label="Role"
                onChange={(e) => setRole(e.target.value)}
              >
                <MenuItem value="user">User</MenuItem>
                <MenuItem value="admin">Admin</MenuItem>
              </Select>
            </FormControl>

            <Button
              type="submit"
              variant="contained"
              size="large"
              fullWidth
              sx={{
                py: 1.2,
                mt: 1,
                background: "linear-gradient(to right, #3b82f6, #8b5cf6, #22c55e)",
                fontWeight: "bold",
                color: "white",
                borderRadius: 2,
                boxShadow: "0 4px 20px rgba(0,0,0,0.2)",
                transition: "0.3s",
                "&:hover": {
                  transform: "scale(1.05)",
                  boxShadow: "0 6px 25px rgba(0,0,0,0.3)",
                },
              }}
            >
              Login
            </Button>
          </Box>

          {/* Footer Link */}
          <Typography variant="body2" sx={{ mt: 3 }}>
            Don’t have an account?{" "}
            <a
              href="/register"
              style={{
                color: "#3b82f6",
                fontWeight: 600,
                textDecoration: "none",
              }}
            >
              Register
            </a>
          </Typography>
        </Paper>
      </Container>
    </Box>
  );
};

export default Login;
