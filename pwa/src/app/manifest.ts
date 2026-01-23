import type { MetadataRoute } from "next";

// Required for static export
export const dynamic = "force-static";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "noise.sh Companion",
    short_name: "noise.sh",
    description: "Capture song ideas on the go - voice memos, photos, tap tempo, and quick notes that sync to your desktop.",
    start_url: "/",
    display: "standalone",
    background_color: "#0A0E27",
    theme_color: "#D4A574",
    orientation: "portrait-primary",
    icons: [
      {
        src: "/icons/icon-192.svg",
        sizes: "192x192",
        type: "image/svg+xml",
        purpose: "any",
      },
      {
        src: "/icons/icon-512.svg",
        sizes: "512x512",
        type: "image/svg+xml",
        purpose: "any",
      },
    ],
    categories: ["music", "productivity", "utilities"],
    lang: "en",
    dir: "ltr",
    prefer_related_applications: false,
  };
}
