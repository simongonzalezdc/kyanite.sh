/**
 * Theme tokens extracted from Go TUI themes
 * Source: internal/theme/registry.go
 */

export interface Theme {
  id: string;
  name: string;
  primary: string;
  secondary: string;
  accent: string;
  background: string;
  text: string;
  success: string;
  warning: string;
  error: string;
}

export const themes: Record<string, Theme> = {
  monochrome: {
    id: "monochrome",
    name: "Monochrome",
    primary: "#FFFFFF",
    secondary: "#999999",
    accent: "#FFFFFF",
    background: "#000000",
    text: "#FFFFFF",
    success: "#CCCCCC",
    warning: "#888888",
    error: "#666666",
  },
  "amber-night": {
    id: "amber-night",
    name: "Amber Night",
    primary: "#D4A574",
    secondary: "#9D84B7",
    accent: "#F4D03F",
    background: "#0A0E27",
    text: "#E8DFF5",
    success: "#52D3AA",
    warning: "#FFA502",
    error: "#EA2027",
  },
  "twilight-mist": {
    id: "twilight-mist",
    name: "Twilight Mist",
    primary: "#B8A3C9",
    secondary: "#8E7B9D",
    accent: "#D4C5E0",
    background: "#151520",
    text: "#E8E4F0",
    success: "#90C695",
    warning: "#C9A87C",
    error: "#C77777",
  },
  "indigo-depths": {
    id: "indigo-depths",
    name: "Indigo Depths",
    primary: "#4169E1",
    secondary: "#5F9EA0",
    accent: "#87CEEB",
    background: "#0A0A1A",
    text: "#E6F2FF",
    success: "#52D3AA",
    warning: "#FFB84D",
    error: "#FF6B6B",
  },
  "forest-path": {
    id: "forest-path",
    name: "Forest Path",
    primary: "#8FBC8F",
    secondary: "#6B8E6B",
    accent: "#B4D7B4",
    background: "#1A1F1A",
    text: "#E8F5E8",
    success: "#90EE90",
    warning: "#DAA520",
    error: "#CD5C5C",
  },
  "clay-earth": {
    id: "clay-earth",
    name: "Clay Earth",
    primary: "#A0522D",
    secondary: "#8B4513",
    accent: "#DEB887",
    background: "#1A1410",
    text: "#F5E6D3",
    success: "#8FBC8F",
    warning: "#CD853F",
    error: "#CD5C5C",
  },
  "iron-forge": {
    id: "iron-forge",
    name: "Iron Forge",
    primary: "#DC143C",
    secondary: "#4A4A4A",
    accent: "#FF6347",
    background: "#1A0A0A",
    text: "#FFE6E6",
    success: "#90C695",
    warning: "#FFB84D",
    error: "#FF4444",
  },
  sunlight: {
    id: "sunlight",
    name: "Sunlight",
    primary: "#FFD700",
    secondary: "#DAA520",
    accent: "#FFF8DC",
    background: "#0F0F0A",
    text: "#FFFACD",
    success: "#98D982",
    warning: "#FF9800",
    error: "#D32F2F",
  },
  "cyan-wave": {
    id: "cyan-wave",
    name: "Cyan Wave",
    primary: "#00CED1",
    secondary: "#4682B4",
    accent: "#7FFFD4",
    background: "#0A1418",
    text: "#E0F7FA",
    success: "#52D3AA",
    warning: "#FFB84D",
    error: "#FF6B6B",
  },
  "electric-rose": {
    id: "electric-rose",
    name: "Electric Rose",
    primary: "#FF1493",
    secondary: "#C71585",
    accent: "#00CED1",
    background: "#1A0A1A",
    text: "#FFF0F8",
    success: "#52D3AA",
    warning: "#FFB84D",
    error: "#FF4444",
  },
};

export const defaultThemeId = "amber-night";

export const getTheme = (id: string): Theme => {
  return themes[id] || themes[defaultThemeId];
};

export const listThemes = (): string[] => [
  "monochrome",
  "amber-night",
  "twilight-mist",
  "indigo-depths",
  "forest-path",
  "clay-earth",
  "iron-forge",
  "sunlight",
  "cyan-wave",
  "electric-rose",
];
