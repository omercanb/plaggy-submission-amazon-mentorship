"use client";

import * as React from "react";
import {
  Stack,
  Breadcrumbs,
  Typography,
  Link as MUILink,
  Tooltip,
  IconButton,
} from "@mui/material";

export type Crumb = {
  label: string;
  href?: string; // if provided, rendered as a link
  onClick?: () => void; // overrides href if present
};

export type HeaderAction = {
  title: string; // tooltip text
  icon: React.ReactNode; // <ArrowBackIcon /> etc.
  onClick: () => void;
  disabled?: boolean;
};

type PageHeaderProps = {
  breadcrumbs?: Crumb[];
  actions?: HeaderAction[]; // right-side icon buttons
};

export default function PageHeader({
  breadcrumbs = [],
  actions,
}: PageHeaderProps) {
  return (
    <Stack direction="column" spacing={1.5} sx={{ mb: 2 }}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <Breadcrumbs aria-label="breadcrumb" sx={{ color: "text.secondary" }}>
          {breadcrumbs.map((c, i) =>
            i < breadcrumbs.length - 1 ? (
              <MUILink
                key={i}
                component="button"
                color="inherit"
                underline="hover"
                onClick={c.onClick}
                {...(c.href && !c.onClick ? { href: c.href } : {})}
              >
                {c.label}
              </MUILink>
            ) : (
              <Typography key={i} color="text.primary">
                {c.label}
              </Typography>
            )
          )}
        </Breadcrumbs>

        {actions && actions.length > 0 && (
          <Stack direction="row" spacing={1}>
            {actions.map((a, i) => (
              <Tooltip key={i} title={a.title}>
                <span>
                  <IconButton
                    size="small"
                    onClick={a.onClick}
                    disabled={a.disabled}
                  >
                    {a.icon}
                  </IconButton>
                </span>
              </Tooltip>
            ))}
          </Stack>
        )}
      </Stack>
    </Stack>
  );
}
