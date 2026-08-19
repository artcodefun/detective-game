import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "ДетектИИв — интерактивная детективная игра",
  description: "Допрашивайте подозреваемых, изучайте улики и раскройте преступление, созданное специально для вас.",
  metadataBase: new URL("https://detective-game.artcodefun.com"),
  openGraph: {
    title: "ДетектИИв — интерактивная детективная игра",
    description: "Допрашивайте подозреваемых, изучайте улики и раскройте преступление.",
    images: [{ url: "/og.png", width: 1200, height: 630, alt: "ДетектИИв" }],
    locale: "ru_RU",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "ДетектИИв — интерактивная детективная игра",
    description: "Допрашивайте подозреваемых, изучайте улики и раскройте преступление.",
    images: ["/og.png"],
  },
  icons: {
    icon: [
      { url: "/favicon.png", type: "image/png", sizes: "32x32" },
      { url: "/icon-512.png", type: "image/png", sizes: "512x512" },
    ],
    shortcut: "/favicon.png",
    apple: [{ url: "/apple-touch-icon.png", sizes: "180x180", type: "image/png" }],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ru">
      <body>{children}</body>
    </html>
  );
}
