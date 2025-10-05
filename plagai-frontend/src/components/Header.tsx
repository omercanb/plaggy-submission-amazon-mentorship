"use client";

import React from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthProvider";
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
  Chip,
  Container,
  useScrollTrigger,
  alpha,
  Button,
} from "@mui/material";
import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import AssignmentTurnedInIcon from "@mui/icons-material/AssignmentTurnedIn";
import AssignHomeworkDialog from "@/components/AssignHomeworkDialog";

function initialsFromEmail(email?: string | null) {
  if (!email) return "U";
  const name = email
    .split("@")[0]
    .replace(/[._-]+/g, " ")
    .trim();
  const parts = name.split(" ").filter(Boolean);
  const first = parts[0]?.[0] ?? "U";
  const second = parts[1]?.[0] ?? "";
  return (first + second).toUpperCase();
}

const Header = () => {
  const [assignOpen, setAssignOpen] = React.useState(false);
  const { user, logout } = useAuth();
  const router = useRouter();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const trigger = useScrollTrigger({ disableHysteresis: true, threshold: 2 });

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };
  const handleMenuClose = () => setAnchorEl(null);

  const handleLogout = () => {
    handleMenuClose();
    logout();
  };

  const userInitials = initialsFromEmail(user?.email);

  return (
    <AppBar
      position="sticky"
      elevation={trigger ? 3 : 0}
      sx={{
        backgroundColor: "rgba(230, 240, 250, 0.85)",
        backdropFilter: "blur(10px)",
        borderBottom: (theme) => `1px solid ${theme.palette.divider}`,
        boxShadow: trigger ? "0 2px 8px rgba(100, 150, 200, 0.25)" : "none",
        color: "text.primary",
        transition: "all 0.2s ease",
      }}
    >
      <Container maxWidth="lg">
        <Toolbar
          disableGutters
          sx={{
            minHeight: 72,
            gap: 1.5,
          }}
        >
          <Box
            onClick={() => router.push("/")}
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 1.25,
              cursor: "pointer",
              px: 1.25,
              py: 0.75,
              borderRadius: 2,
              transition: "background-color 160ms ease",
              "&:hover": (t) => ({
                backgroundColor: alpha(t.palette.primary.main, 0.06),
              }),
              position: "relative",
              overflow: "hidden",
              "&::after": {
                content: '""',
                position: "absolute",
                left: 10,
                right: 10,
                bottom: 0,
                height: 2,
                background:
                  "linear-gradient(90deg, rgba(99,102,241,1) 0%, rgba(59,130,246,1) 100%)",
                borderRadius: 1,
                opacity: 0.6,
              },
            }}
          >
            <Box
              sx={{
                width: 28,
                height: 28,
                borderRadius: "8px",
                background:
                  "linear-gradient(135deg, rgba(99, 132, 204, 1) 0%, rgba(150, 180, 230, 1) 100%)",
                boxShadow: "0 2px 6px rgba(99, 132, 204, 0.35)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                width="18"
                height="18"
                fill="white"
              >
                <path d="M7 4h6a5 5 0 010 10H9v6H7V4zm2 8h4a3 3 0 000-6H9v6z" />
              </svg>
            </Box>
            <Typography
              variant="h6"
              sx={{
                fontWeight: 800,
                letterSpacing: 0.2,
                color: "text.primary",
              }}
            >
              Plaggy
            </Typography>
          </Box>
          <Box sx={{ flex: 1 }} />
          {user && (
            <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
              <Button
                variant="contained"
                size="small"
                startIcon={<AssignmentTurnedInIcon />}
                onClick={() => setAssignOpen(true)}
                sx={{ mr: 1, display: { xs: "none", sm: "inline-flex" } }}
              >
                Assign homework
              </Button>
              <Chip
                label={user.email}
                size="small"
                variant="outlined"
                sx={{
                  display: { xs: "none", sm: "inline-flex" },
                  borderRadius: 2,
                }}
              />
              <IconButton
                size="small"
                onClick={handleMenuOpen}
                aria-controls="user-menu"
                aria-haspopup="true"
                aria-label="account menu"
                sx={{
                  "&:hover": {
                    backgroundColor: "rgba(99,102,241,0.08)",
                  },
                }}
              >
                <Avatar
                  sx={{
                    width: 34,
                    height: 34,
                    fontSize: 14,
                    background:
                      "linear-gradient(135deg, rgba(99, 132, 204, 1) 0%, rgba(150, 180, 230, 1) 100%)",
                    boxShadow: "0 2px 6px rgba(99, 132, 204, 0.35)",
                    color: "#fff",
                  }}
                >
                  {userInitials}
                </Avatar>
              </IconButton>
              <Menu
                id="user-menu"
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
                anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
                transformOrigin={{ vertical: "top", horizontal: "right" }}
              >
                <MenuItem onClick={handleLogout}>
                  <LogoutRoundedIcon
                    fontSize="small"
                    style={{ marginRight: 8 }}
                  />
                  Logout
                </MenuItem>
              </Menu>
            </Box>
          )}
        </Toolbar>
      </Container>
      <AssignHomeworkDialog
        open={assignOpen}
        onClose={() => setAssignOpen(false)}
      />
    </AppBar>
  );
};

export default Header;
