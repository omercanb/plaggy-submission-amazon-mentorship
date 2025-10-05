"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "../providers/AuthProvider";
import { CircularProgress, Box } from "@mui/material";

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth(); // custom hook
  const router = useRouter();

  useEffect(() => {
    // If we are not loading and we don't have a user object,
    // this probably means we fucked up somewhere and should login again ~brtcrt
    if (!loading && !user) {
      router.push("/login");
    }
  }, [user, loading, router]);

  if (loading) {
    // visual feedback. yay. fucking end my suffering ~brtcrt
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  if (!user) {
    return null; // Redirect will happen in useEffect, so no need to router.push here
  }

  return children;
}
