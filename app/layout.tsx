import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";

const inter = localFont({
  src: [
    {
      path: "./fonts/Inter-VariableFont_opsz,wght.ttf",
      weight: "100 900",
      style: "normal",
    },
    {
      path: "./fonts/Inter-Italic-VariableFont_opsz,wght.ttf",
      weight: "100 900",
      style: "italic",
    },
  ],
  variable: "--font-inter",
  display: "swap",
});

const spaceMono = localFont({
  src: [
    {
      path: "./fonts/SpaceMono-Regular.ttf",
      weight: "400",
      style: "normal",
    },
    {
      path: "./fonts/SpaceMono-Italic.ttf",
      weight: "400",
      style: "italic",
    },
    {
      path: "./fonts/SpaceMono-Bold.ttf",
      weight: "700",
      style: "normal",
    },
    {
      path: "./fonts/SpaceMono-BoldItalic.ttf",
      weight: "700",
      style: "italic",
    },
  ],
  variable: "--font-space-mono",
  display: "swap",
});

const outfit = localFont({
  src: [
    {
      path: "./fonts/Outfit-VariableFont_wght.ttf",
      weight: "100 900",
      style: "normal",
    },
  ],
  variable: "--font-outfit",
  display: "swap",
});

export const metadata: Metadata = {
  title: "CodeWave Academy",
  description: "Interactive coding platform with AI-generated exercises",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${spaceMono.variable} ${outfit.variable}`}>
      <head>
        <link rel="icon" href="/favicon.ico" />
      </head>
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
