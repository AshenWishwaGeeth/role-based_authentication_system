import React, { useState } from "react";
import {
  Box,
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  MenuItem,
} from "@mui/material";

const Register: React.FC = () => {
  const [name, setName] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [role, setRole] = useState<string>("user");

  return (
    <Box
      sx={{
        minHeight: "100vh",
        background: "linear-gradient(to bottom right, #bbf7d0, #c7d2fe, #a5b4fc)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Container maxWidth="xs">
        <Paper
          elevation={8}
          sx={{
            p: 4,
            borderRadius: 4,
            textAlign: "center",
            backdropFilter: "blur(10px)",
          }}
        >
          <Typography variant="h4" fontWeight="bold" color="primary" mb={2}>
            Create Account
          </Typography>
          <Typography variant="body1" color="text.secondary" mb={3}>
            Sign up to get started
          </Typography>

          <Box
            component="form"
            sx={{ display: "flex", flexDirection: "column", gap: 3 }}
          >
            <TextField
              label="Name"
              variant="outlined"
              fullWidth
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
            <TextField
              label="Email"
              type="email"
              variant="outlined"
              fullWidth
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <TextField
              label="Role"
              select
              variant="outlined"
              fullWidth
              value={role}
              onChange={(e) => setRole(e.target.value)}
              required
            >
              <MenuItem value="user">User</MenuItem>
              <MenuItem value="admin">Admin</MenuItem>
            </TextField>
            <TextField
              label="Password"
              type="password"
              variant="outlined"
              fullWidth
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />

            <Button
              type="submit"
              variant="contained"
              size="large"
              sx={{
                py: 1.2,
                mt: 1,
                background: "linear-gradient(to right, #22c55e, #3b82f6, #8b5cf6)",
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
              Register
            </Button>
          </Box>

          <Typography variant="body2" sx={{ mt: 3 }}>
            Already have an account?{" "}
            <a
              href="/login"
              style={{
                color: "#3b82f6",
                fontWeight: 600,
                textDecoration: "none",
              }}
            >
              Login
            </a>
          </Typography>
        </Paper>
      </Container>
    </Box>
  );
};

export default Register;
