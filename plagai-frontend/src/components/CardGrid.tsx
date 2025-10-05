"use client";

import * as React from "react";
import {
  Grid,
  Card,
  CardActionArea,
  CardContent,
  Typography,
  Chip,
  Box,
  Skeleton,
} from "@mui/material";
import { alpha } from "@mui/material/styles";
import { useRouter } from "next/navigation";

export type CardsGridItem = {
  id: string;
  title: string;
  subtitle?: string; // e.g., ID, due date, small note
  href?: string; // if provided, clicking navigates
  onClick?: () => void; // overrides href if given
  rightLabel?: string; // small chip on the right (optional)
};

type CardsGridProps = {
  items: CardsGridItem[];
  loading?: boolean;
  emptyText?: string;
  columns?: { xs?: number; sm?: number; md?: number; lg?: number }; // per-breakpoint
  cardMinHeight?: number;
};

export default function CardsGrid({
  items,
  loading = false,
  emptyText = "Nothing to show yet",
  cardMinHeight = 110,
}: CardsGridProps) {
  const router = useRouter();

  if (loading) {
    return (
      <Grid container spacing={3}>
        {Array.from({ length: 6 }).map((_, i) => (
          <Grid key={i}>
            <Card
              variant="outlined"
              sx={{
                borderRadius: 3,
                overflow: "hidden",
                borderColor: (t) => alpha(t.palette.primary.main, 0.15),
                backgroundColor: (t) => alpha(t.palette.primary.light, 0.04),
              }}
            >
              <Skeleton variant="rectangular" height={cardMinHeight} />
              <Box sx={{ p: 2 }}>
                <Skeleton variant="text" width="60%" />
                <Skeleton variant="text" width="40%" />
              </Box>
            </Card>
          </Grid>
        ))}
      </Grid>
    );
  }

  if (!items?.length) {
    return (
      <Box
        sx={{
          py: 8,
          textAlign: "center",
          borderRadius: 3,
          border: (t) => `1px dashed ${t.palette.divider}`,
          backgroundColor: (t) => alpha(t.palette.primary.light, 0.03),
        }}
      >
        <Typography variant="h6" gutterBottom>
          {emptyText}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Check back later.
        </Typography>
      </Box>
    );
  }

  return (
    <Grid container spacing={3}>
      {items.map((it) => {
        const handleClick = () => {
          if (it.onClick) return it.onClick();
          if (it.href) router.push(it.href);
        };

        return (
          <Grid key={it.id}>
            <Card
              elevation={3}
              sx={{
                borderRadius: 3,
                overflow: "hidden",
                transition: "transform 180ms ease, box-shadow 180ms ease",
                "&:hover": { transform: "translateY(-2px)", boxShadow: 6 },
                "&::before": {
                  content: '""',
                  display: "block",
                  height: 6,
                  background:
                    "linear-gradient(135deg, rgba(99, 132, 204, 1) 0%, rgba(150, 180, 230, 1) 100%)",
                },
              }}
            >
              <CardActionArea
                onClick={handleClick}
                sx={{ minHeight: cardMinHeight }}
              >
                <CardContent sx={{ p: 2.5 }}>
                  <Box
                    display="flex"
                    alignItems="center"
                    justifyContent="space-between"
                    gap={1}
                  >
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                      {it.title}
                    </Typography>
                    {it.rightLabel && (
                      <Chip
                        label={it.rightLabel}
                        size="small"
                        variant="outlined"
                        sx={{
                          borderRadius: 2,
                          borderColor: (t) =>
                            alpha(t.palette.primary.main, 0.35),
                          bgcolor: (t) => alpha(t.palette.primary.light, 0.08),
                        }}
                      />
                    )}
                  </Box>

                  {it.subtitle && (
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      sx={{ mt: 0.75 }}
                    >
                      {it.subtitle}
                    </Typography>
                  )}
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        );
      })}
    </Grid>
  );
}
