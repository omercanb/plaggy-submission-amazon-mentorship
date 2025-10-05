"use client";
import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  palette: {
    primary: {
      main: "#6384cc",
      light: "#96b4e6",
      dark: "#4a6bab",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#3aa8b4",
      light: "#6fd3da",
      dark: "#2a7d85",
      contrastText: "#ffffff",
    },
    background: {
      default: "#f5f8fc",
      paper: "#ffffff",
    },
  },
});

export default theme;
